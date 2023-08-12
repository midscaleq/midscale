package group

import (
	"encoding/json"
	"fmt"
	groupModel "midscale/midscale/app/data/model/group"
	"midscale/midscale/app/data/model/httpbase"
	"midscale/midscale/app/service/group"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Join Group godoc
// @Summary Join a group.
// @Description
// @Description
// @ID groupJoin
// @Tags Group
// @Accept  json
// @Produce  json
// @Success 200 {object} "code=0 成功;"
// @Router /ms/v1/api/group/join [POST]
func Join(c *gin.Context) {
	var encodedReq groupModel.EncodedReq
	if err := c.ShouldBindJSON(&encodedReq); err != nil {
		c.JSON(http.StatusOK, httpbase.NewBaseRsp(c, httpbase.ParmsError,
			fmt.Errorf("failed BindJSON %v", err).Error(), nil))
		return
	}
	data, sign, err := group.UnpackConnectionInfo(encodedReq.EncodedReq)
	if err != nil {
		c.JSON(http.StatusOK, httpbase.NewBaseRsp(c, httpbase.ParmsError,
			fmt.Errorf("failed UnpackConnectionInfo %v", err).Error(), nil))
		return
	}

	var joinReq groupModel.JoinReq
	err = json.Unmarshal(data, &joinReq)
	if err != nil {
		c.JSON(http.StatusOK, httpbase.NewBaseRsp(c, httpbase.ParmsError,
			fmt.Errorf("failed Unmarshal %v", err).Error(), nil))
		return
	}

	veryfied, err := group.VerifySign(joinReq.ServerTransportPubKey, data, sign)
	if err != nil {
		c.JSON(http.StatusOK, httpbase.NewBaseRsp(c, httpbase.ParmsError,
			fmt.Errorf("failed VerifySign %v", err).Error(), nil))
		return
	}
	if !veryfied {
		c.JSON(http.StatusOK, httpbase.NewBaseRsp(c, httpbase.ParmsError, "sign invalid", nil))
		return
	}

	rsp, err := group.NewJoinService().Join(encodedReq.JoinPassword, encodedReq.JoinerNickName,
		joinReq.ServerTransportPubKey, joinReq.Host, joinReq.GroupName, joinReq.ServerNickName,
		encodedReq.JoinerListenPort, encodedReq.ShowConfQR)
	if err != nil {
		c.JSON(http.StatusOK, httpbase.NewBaseRsp(c, httpbase.ServerError, err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, httpbase.NewOkBaseRsp(c, rsp))
}
