package group

import (
	"fmt"
	groupinfo "midscale/midscale/app/data/mgr/group/info"
	"midscale/midscale/app/data/mgr/wg/conf"
	"midscale/midscale/app/data/model/define"
	groupModel "midscale/midscale/app/data/model/group"
	modelinfo "midscale/midscale/app/data/model/group/info"
	"midscale/midscale/deps/mskit/net/utils"
)

type ApplyService struct {
}

func NewApplyService() *ApplyService {
	return &ApplyService{}
}

// return serverTransportPubKey, serverNickName, joinerLocalIP, error
func (applyService *ApplyService) Apply(joinerTransportPubKey, groupName, serverNickName,
	joinPassword, joinerNickName, joinerPubKey string,
	joinerListenPort int) (*groupModel.ApplyRsp, error) {

	verified, err := groupinfo.VerifyJoinPassword(groupName, serverNickName, joinPassword)
	if err != nil {
		return nil, fmt.Errorf("group not exist")
	}
	if !verified {
		return nil, fmt.Errorf("joining password is invalid")
	}

	exist, err := conf.ConfExist(groupName, serverNickName)
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, fmt.Errorf("group of conf not exist")
	}

	joiner := &modelinfo.JoinerInfo{JoinerTransportPubKey: joinerTransportPubKey,
		JoinerNickName:   joinerNickName,
		JoinerPubKey:     joinerPubKey,
		JoinerLocalIP:    "", // want to auto dispatch a new joiner ip.
		JoinerListenPort: joinerListenPort}
	err = groupinfo.AppendJoiner(groupName, serverNickName, joiner)
	if err != nil {
		return nil, err
	}
	groupInfo, err := groupinfo.GetGroupInfo(groupName, serverNickName)
	if err != nil {
		return nil, err
	}

	replacedJoinerLocalIP, err := utils.ReplaceCIDRStringWithNewNetMask(joiner.JoinerLocalIP, "32")
	if err != nil {
		return nil, err
	}
	allowedIPs := []string{replacedJoinerLocalIP}
	err = conf.AppendPeerStr(groupName, serverNickName, "", joinerPubKey, allowedIPs, define.PersistentKeepalive)
	if err != nil {
		return nil, err
	}

	wgPeers, err := conf.GetPeers(groupName, serverNickName,
		fmt.Sprintf("%v:%v", define.EndpointHost, define.DefaultListenPort),
		define.PersistentKeepalive)
	if err != nil {
		return nil, err
	}

	groupInfo.ServerJoinPassword = ""
	groupInfo.AutoStartTunnel = false
	rsp := &groupModel.ApplyRsp{Accepted: true,
		GroupInfo:     groupInfo,
		JoinerLocalIP: joiner.JoinerLocalIP,
		WGPeers:       wgPeers}

	if define.AutoStartTunnel {
		if err := StopTunnel(groupName); err != nil {
			return nil, fmt.Errorf("JoinService.Apply, StopTunnel err: %v", err)
		}
		if err := StartTunnel(groupName, serverNickName); err != nil {
			return nil, fmt.Errorf("JoinService.Apply, StartTunnel err: %v", err)
		}
	}

	return rsp, nil
}
