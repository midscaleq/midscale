package group

import (
	"midscale/midscale/app/data/model/group/info"
	"midscale/midscale/deps/wgconf"
)

type ApplyReq struct {
	JoinerTransportPubKey string `json:"joinerTransportPubKey" binding:"required"`
	GroupName             string `json:"groupName" binding:"required"`
	ServerNickName        string `json:"serverNickName" binding:"required"`
	JoinPassword          string `json:"joinPassword"`
	JoinerNickName        string `json:"joinerNickName" binding:"required"`
	JoinerPubKey          string `json:"joinerPubKey" binding:"required"`
	JoinerListenPort      int    `json:"joinerListenPort"`
}

type ApplyRsp struct {
	Accepted      bool            `json:"accepted"`
	JoinerLocalIP string          `json:"joinerLocalIP"`
	GroupInfo     *info.GroupInfo `json:"groupInfo"`
	WGPeers       []*wgconf.Peer  `json:"members"`
}
