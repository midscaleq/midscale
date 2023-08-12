package group

import (
	"fmt"
	"midscale/midscale/app/data/mgr/transport/key"
	groupModel "midscale/midscale/app/data/model/group"
	"midscale/midscale/app/data/model/httpbase"
	"midscale/midscale/app/service/group"
	"midscale/midscale/deps/mskit/encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Apply Group godoc
// @Summary Apply a group.
// @Description curl.exe -X POST http://localhost:18888/ms/v1/api/group/apply -H 'Content-Type: application/json' -d '{ \"groupName\": \"ms-0704\", \"clientPubKey\":\"\" }'
// @ID groupApply
// @Tags Group
// @Accept  json
// @Produce  json
// @Success 200 {object} "code=0 成功;"
// @Router /ms/v1/api/group/apply [POST]
func Apply(c *gin.Context) {
	var encReq httpbase.EncryptedReq
	if err := c.ShouldBindJSON(&encReq); err != nil {
		c.JSON(http.StatusOK, httpbase.NewBaseRsp(c, httpbase.ParmsError,
			fmt.Errorf("failed BindJSON %v", err).Error(), nil))
		return
	}

	// decrypt data by ServerTransportPriKey
	var req groupModel.ApplyReq
	if err := json.DecryptB64ToJson(key.GetTransportPriKeyHex(), encReq.EncryptedReq, encReq.PlainHash, &req); err != nil {
		c.JSON(http.StatusOK, httpbase.NewBaseRsp(c, httpbase.ParmsError,
			fmt.Errorf("failed DecryptB64ToJson %v", err).Error(), nil))
		return
	}

	rsp, err := group.NewApplyService().Apply(req.JoinerTransportPubKey, req.GroupName, req.ServerNickName,
		req.JoinPassword, req.JoinerNickName, req.JoinerPubKey, req.JoinerListenPort)
	if err != nil {
		c.JSON(http.StatusOK, httpbase.NewBaseRsp(c, httpbase.ServerError, err.Error(), nil))
		return
	}

	// encrypt data by JoinerTransportPubKey
	rspB64, plainHash, err := json.EncryptJsonToB64(req.JoinerTransportPubKey, rsp)
	if err != nil {
		c.JSON(http.StatusOK, httpbase.NewBaseRsp(c, httpbase.ParmsError, err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, httpbase.NewOkBaseRsp(c, &httpbase.EncryptedRsp{EncryptedRsp: rspB64, PlainHash: plainHash}))
}
