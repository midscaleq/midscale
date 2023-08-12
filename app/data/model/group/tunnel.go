package group

const (
	OperationStart = "start"
	OperationStop  = "stop"
)

type TunnelReq struct {
	GroupName string `json:"groupName" binding:"required"`
	NickName  string `json:"nickName"`
	Operation string `json:"operation" binding:"required"`
}
