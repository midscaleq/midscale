package info

type GroupInfoReq struct {
	GroupName string `json:"groupName" binding:"required"`
	NickName  string `json:"nickName" binding:"required"`
}

type GroupInfo struct {
	GroupName             string `json:"groupName"`
	ServerJoinPassword    string `json:"serverjoinPassword,omitempty"`
	ServerNickName        string `json:"serverNickName,omitempty"`
	ServerPubKey          string `json:"serverPubKey"`
	ServerLocalIP         string `json:"serverLocalIP"`
	ServerLocalIPCIDR     string `json:"serverLocalIPCIDR"`
	ServerListenPort      int    `json:"serverListenPort"`
	EncodedConnectionInfo string `json:"encodedConnectionInfo"`
	ServerTransportPubKey string `json:"serverTransportPubKey"`
	AutoStartTunnel       bool   `json:"autoStartTunnel"`

	JoinerInfos []*JoinerInfo `json:"joinerInfos"`
}

type JoinerInfo struct {
	JoinerTransportPubKey string `json:"joinerTransportPubKey"`
	JoinerNickName        string `json:"joinerNickName,omitempty"`
	JoinerPubKey          string `json:"joinerPubKey"`
	JoinerLocalIP         string `json:"joinerLocalIP"`
	JoinerListenPort      int    `json:"joinerListenPort"`
}
