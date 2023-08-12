package group

import (
	"fmt"
	groupinfo "midscale/midscale/app/data/mgr/group/info"
	"midscale/midscale/app/data/mgr/transport/key"
	"midscale/midscale/app/data/mgr/wg/conf"
	"midscale/midscale/app/data/model/define"
	groupModel "midscale/midscale/app/data/model/group"
	modelinfo "midscale/midscale/app/data/model/group/info"
	"time"
)

type CreateService struct {
}

func NewCreateService() *CreateService {
	return &CreateService{}
}

func (createService *CreateService) Create(groupName, joinPassword, nickName, localIP string,
	oneMaskLength, listenPort int,
	DNS string) (string, error) {

	if len(localIP) < 1 {
		localIP = define.DefaultLocalIP
	}
	if oneMaskLength < 1 {
		oneMaskLength = define.DefaultOneMaskLength
	}
	ServerLocalIPCIDR := fmt.Sprintf("%v/%v", localIP, oneMaskLength)
	if listenPort < 1 {
		listenPort = define.DefaultListenPort
	}
	if len(DNS) < 1 {
		DNS = define.DefaultDNS
	}

	exist, err := groupinfo.GroupExist(groupName)
	if err != nil {
		return "", fmt.Errorf("CreateService.Create, GroupExist err: %v", err)
	}
	if exist {
		return "", fmt.Errorf("CreateService.Create, GetGroupInfo err, group already exist")
	}

	err = conf.NewConf(groupName, nickName, localIP,
		oneMaskLength, listenPort,
		DNS)
	if err != nil {
		return "", fmt.Errorf("CreateService.Create, NewConf err: %v", err)
	}
	pubKey, err := conf.GetPubKey(groupName, nickName)
	if err != nil {
		return "", fmt.Errorf("CreateService.Create, GetPubKey err: %v", err)
	}

	now := time.Now().In(time.UTC)
	expiredAt := now.Add(24 * time.Hour * 365 * 3)
	host := fmt.Sprintf("%v%v", define.ServerHost, define.DefaultListenAddrForRemote)
	connInfo := &groupModel.ConnectionInfo{ServerTransportPubKey: key.GetTransportPubKeyHex(),
		Host:           host,
		GroupName:      groupName,
		ServerNickName: nickName,
		Iat:            now.String(),
		Exp:            expiredAt.String()}
	encodedConnectionInfo, err := PackConnectionInfo(connInfo)
	if err != nil {
		return "", fmt.Errorf("CreateService.Create, PackConnectionInfo err: %v", err)
	}

	groupInfo := &modelinfo.GroupInfo{GroupName: groupName,
		ServerJoinPassword:    joinPassword,
		ServerNickName:        nickName,
		ServerPubKey:          pubKey,
		ServerLocalIP:         localIP,
		ServerLocalIPCIDR:     ServerLocalIPCIDR,
		ServerListenPort:      listenPort,
		EncodedConnectionInfo: encodedConnectionInfo,
		ServerTransportPubKey: key.GetTransportPubKeyHex(),
		AutoStartTunnel:       define.AutoStartTunnel}
	err = groupinfo.NewGroupInfo(groupInfo)
	if err != nil {
		return "", fmt.Errorf("CreateService.Create, NewGroupInfo err: %v", err)
	}

	return groupInfo.EncodedConnectionInfo, nil
}
