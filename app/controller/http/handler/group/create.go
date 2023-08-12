package group

import (
	"fmt"
	"net/http"

	groupModel "midscale/midscale/app/data/model/group"
	"midscale/midscale/app/data/model/httpbase"
	"midscale/midscale/app/service/group"

	"github.com/gin-gonic/gin"
)

// Create Group godoc
// @Summary Create a group.
// @Description
// @ID groupCreate
// @Tags Group
// @Accept  json
// @Produce  json
// @Success 200 {object} "code=0 成功;"
// @Router /ms/v1/api/group/create [POST]
func Create(c *gin.Context) {

	var req groupModel.CreateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, httpbase.NewBaseRsp(c, httpbase.ParmsError,
			fmt.Errorf("failed BindJSON %v", err).Error(), nil))
		return
	}

	encodedConnectionInfo, err := group.NewCreateService().Create(req.GroupName,
		req.JoinPassword,
		req.NickName,
		req.LocalIP,
		req.OneMaskLength,
		req.ListenPort,
		req.DNS)
	if err != nil {
		c.JSON(http.StatusOK, httpbase.NewBaseRsp(c, httpbase.ServerError, err.Error(), nil))
		return
	}

	rsp := &groupModel.CreateRsp{EncodedConnectionInfo: encodedConnectionInfo}
	c.JSON(http.StatusOK, httpbase.NewOkBaseRsp(c, rsp))
}
