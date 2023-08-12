package group

type CreateReq struct {
	GroupName     string `json:"groupName" binding:"required"`
	JoinPassword  string `json:"joinPassword"`
	NickName      string `json:"nickName" binding:"required"`
	LocalIP       string `json:"localIP"`
	OneMaskLength int    `json:"oneMaskLength"`
	ListenPort    int    `json:"listenPort"`
	DNS           string `json:"dNS"`
}

type CreateRsp struct {
	EncodedConnectionInfo string `json:"encodedConnectionInfo"`
}

type ConnectionInfo struct {
	ServerTransportPubKey string `json:"serverTransportPubKey"`
	Host                  string `json:"host"`
	GroupName             string `json:"groupName"`
	ServerNickName        string `json:"serverNickName"`
	Iat                   string `json:"createdTm"`
	Exp                   string `json:"expiredTm"`
}
