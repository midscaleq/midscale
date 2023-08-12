package group

import (
	"fmt"
	groupinfo "midscale/midscale/app/data/mgr/group/info"
	"midscale/midscale/app/data/mgr/transport/key"
	"midscale/midscale/app/data/mgr/wg/conf"
	"midscale/midscale/app/data/model/define"
	groupModel "midscale/midscale/app/data/model/group"
	"midscale/midscale/app/data/model/httpbase"
	"midscale/midscale/deps/mskit/encoding/json"
	"midscale/midscale/deps/mskit/net/http"

	uuid "github.com/satori/go.uuid"
)

const (
	PlaceholderLocalIP = "255.255.255.255"
	ApplyPath          = "/ms/v1/api/group/apply"
)

type JoinService struct {
}

func NewJoinService() *JoinService {
	return &JoinService{}
}

func getJoinerWgConf(groupName, joinerNickName, localIP, dns string, joinerListenPort int) (string, error) {
	err := conf.NewConf(groupName, joinerNickName, localIP, 0, joinerListenPort, dns)
	if err != nil {
		return "", err
	}
	return conf.GetPubKey(groupName, joinerNickName)
}

func (joinService *JoinService) Join(joinPassword, joinerNickName,
	serverTransportPubKey, host, groupName, serverNickName string,
	joinerListenPort int, showConfQR bool) (*groupModel.JoinRsp, error) {
	if joinerListenPort < 1 {
		joinerListenPort = define.DefaultListenPort
	}

	joined := groupinfo.GroupJoined(groupName, joinerNickName)
	if joined {
		return nil, fmt.Errorf("JoinService.Join, GroupJoined already joined")
	}

	joinerPubKey, err := getJoinerWgConf(groupName, joinerNickName, PlaceholderLocalIP, define.DefaultDNS, joinerListenPort)
	if err != nil {
		return nil, fmt.Errorf("JoinService.Join, getJoinerWgConf err: %v", err)
	}

	req := &groupModel.ApplyReq{JoinerTransportPubKey: key.GetTransportPubKeyHex(),
		GroupName:        groupName,
		ServerNickName:   serverNickName,
		JoinPassword:     joinPassword,
		JoinerNickName:   joinerNickName,
		JoinerPubKey:     joinerPubKey,
		JoinerListenPort: joinerListenPort}

	headers := map[string]string{"Request-Id": uuid.NewV4().String()}

	type CustomBaseRsp struct {
		httpbase.CommonRsp
		Results *httpbase.EncryptedRsp `json:"results,omitempty"`
	}

	// encrypt data by ServerTransportPubKey
	reqB64, plainHash, err := json.EncryptJsonToB64(serverTransportPubKey, req)
	if err != nil {
		return nil, fmt.Errorf("JoinService.Join, EncryptJsonToB64 err: %v", err)
	}

	url := host + ApplyPath
	var rsp CustomBaseRsp
	httpReq, httpRsp, _, err := http.PostJsonWithHeaders(url,
		&httpbase.EncryptedReq{EncryptedReq: reqB64, PlainHash: plainHash}, headers, &rsp)
	if err != nil {
		return nil, fmt.Errorf("JoinService.Join, PostJsonWithHeaders err: %v", err)
	}
	if rsp.Code != httpbase.CodeOk {
		return nil, fmt.Errorf("JoinService.Join, PostJsonWithHeaders failed, rsp.Code:%v, rsp.Message:%v. httpReq:[%+v], httpRsp:[%+v]",
			rsp.Code, rsp.Message, httpReq, httpRsp)
	}
	// decrypt data by ClientTransportPriKey
	var applyRsp groupModel.ApplyRsp
	if err := json.DecryptB64ToJson(key.GetTransportPriKeyHex(),
		rsp.Results.EncryptedRsp, rsp.Results.PlainHash, &applyRsp); err != nil {
		return nil, fmt.Errorf("JoinService.Join, DecryptB64ToJson err: %v", err)
	}
	if !applyRsp.Accepted {
		return nil, fmt.Errorf("JoinService.Join, rejected by server, groupName: %v", groupName)
	}

	err = conf.UpdateLocalIp(groupName, joinerNickName, applyRsp.JoinerLocalIP)
	if err != nil {
		return nil, fmt.Errorf("JoinService.Join, UpdateLocalIp err: %v", err)
	}

	for _, wgPeer := range applyRsp.WGPeers {

		if wgPeer.PublicKey == joinerPubKey { // self , ignore
			continue
		}

		err := conf.AppendPeer(groupName, joinerNickName, wgPeer.Endpoint,
			wgPeer.PublicKey, wgPeer.AllowedIPs,
			wgPeer.PersistentKeepalive)
		if err != nil {
			return nil, fmt.Errorf("JoinService.Join, AppendPeer err: %v", err)
		}
	}

	applyRsp.GroupInfo.AutoStartTunnel = define.AutoStartTunnel
	err = groupinfo.NewGroupInfo(applyRsp.GroupInfo)
	if err != nil {
		return nil, fmt.Errorf("JoinService.Create, NewGroupInfo err: %v", err)
	}

	if define.AutoStartTunnel {
		if err := StopTunnel(groupName); err != nil {
			return nil, fmt.Errorf("JoinService.Join, StopTunnel err: %v", err)
		}
		if err := StartTunnel(groupName, joinerNickName); err != nil {
			return nil, fmt.Errorf("JoinService.Join, StartTunnel err: %v", err)
		}
	}

	confQR := make([]byte, 0)
	if showConfQR {
		confQR, err = conf.GetConfBuf(groupName, joinerNickName)
		if err != nil {
			return nil, fmt.Errorf("JoinService.Create, GetConfBuf err: %v", err)
		}
	}

	joinRsp := &groupModel.JoinRsp{ServerNickName: applyRsp.GroupInfo.ServerNickName,
		ServerLocalIP:  applyRsp.GroupInfo.ServerLocalIP,
		JoinerNickName: joinerNickName,
		JoinerLocalIP:  applyRsp.JoinerLocalIP,
		ConfQR:         string(confQR)}

	return joinRsp, nil
}
