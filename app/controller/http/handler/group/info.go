package group

import (
	"fmt"
	"log"
	"net/http"

	"midscale/midscale/app/data/model/group/info"
	"midscale/midscale/app/data/model/httpbase"
	"midscale/midscale/app/service/group"

	"github.com/gin-gonic/gin"
)

// Get Group Info godoc
// @Summary Get a group info.
// @Description curl.exe http://localhost:18888/ms/v1/api/group/info/ms-server
// @ID groupInfo
// @Tags Group
// @Accept  json
// @Produce  json
// @Success 200 {object} "code=0 成功;"
// @Router /ms/v1/api/group/info [GET]
func GetInfo(c *gin.Context) {

	log.Printf("c.Request:%+v", c.Request)
	var req info.GroupInfoReq
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusOK, httpbase.NewBaseRsp(c, httpbase.ParmsError,
			fmt.Errorf("failed BindJSON %v", err).Error(), nil))
		return
	}

	groupInfo, err := group.NewInfoService().GetInfo(req.GroupName, req.NickName)
	if err != nil {
		c.JSON(http.StatusOK, httpbase.NewBaseRsp(c, httpbase.ServerError, err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, httpbase.NewOkBaseRsp(c, groupInfo))
}
