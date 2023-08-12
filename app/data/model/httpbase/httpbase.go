package httpbase

import (
	"github.com/gin-gonic/gin"
)

const (
	CodeOk      = 0
	ParmsError  = 40000
	ServerError = 50000
)

type EncryptedReq struct {
	EncryptedReq string `json:"encryptedReq" validate:"required"`
	PlainHash    string `json:"plainHash,omitempty"`
}
type EncryptedRsp struct {
	EncryptedRsp string `json:"encryptedRsp"`
	PlainHash    string `json:"plainHash,omitempty"`
}

type CommonRsp struct {
	Code    int    `json:"code"`
	Message string `json:"message,omitempty"`
}

type BaseRsp struct {
	CommonRsp
	Results interface{} `json:"results,omitempty"`
}

func NewOkBaseRsp(c *gin.Context, results interface{}) BaseRsp {
	requestId := c.Request.Header.Get("requestId")
	if len(requestId) > 0 {
		c.Header("requestId", string(requestId[0]))
	}
	return BaseRsp{CommonRsp: CommonRsp{Code: CodeOk, Message: "OK"}, Results: results}
}
func NewBaseRsp(c *gin.Context, code int, message string, results interface{}) *BaseRsp {
	requestId := c.Request.Header.Get("requestId")
	if len(requestId) > 0 {
		c.Header("requestId", string(requestId[0]))
	}
	return &BaseRsp{CommonRsp: CommonRsp{Code: code, Message: message}, Results: results}
}
