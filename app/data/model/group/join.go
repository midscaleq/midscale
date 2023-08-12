package group

type EncodedReq struct {
	EncodedReq       string `json:"encodedReq" binding:"required"`
	JoinPassword     string `json:"joinPassword"`
	JoinerNickName   string `json:"joinerNickName" binding:"required"`
	JoinerListenPort int    `json:"joinerListenPort"`
	ShowConfQR       bool   `json:"showConfQR"`
}
type JoinReq struct {
	ServerTransportPubKey string `json:"serverTransportPubKey" binding:"required"`
	Host                  string `json:"host" binding:"required"`
	GroupName             string `json:"groupName" binding:"required"`
	ServerNickName        string `json:"serverNickName"`
}

type JoinRsp struct {
	ServerNickName string `json:"serverNickName"`
	ServerLocalIP  string `json:"serverLocalIP"`
	JoinerNickName string `json:"joinerNickName"`
	JoinerLocalIP  string `json:"joinerLocalIP"`
	ConfQR         string `json:"confQR"`
}
