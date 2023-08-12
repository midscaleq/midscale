package group

import (
	"fmt"
	"midscale/midscale/app/data/model/group"
	"midscale/midscale/app/data/model/httpbase"
	groupService "midscale/midscale/app/service/group"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Group Tunnel godoc
// @Summary  start/stop the group tunnel.
// @Description curl.exe -X POST http://localhost:18888/ms/v1/api/group/tunnel -H 'Content-Type: application/json' -d '{ \"groupName\": \"ms-0704\",\"operation\": \"stop\" }'
// @ID groupTunnel
// @Tags Group
// @Accept  json
// @Produce  json
// @Success 200 {object} "code=0 成功;"
// @Router /ms/v1/api/group/tunnel [POST]
func Tunnel(c *gin.Context) {
	var req group.TunnelReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, httpbase.NewBaseRsp(c, httpbase.ParmsError,
			fmt.Errorf("failed BindJSON %v", err).Error(), nil))
		return
	}

	if req.Operation == group.OperationStart {
		if len(req.NickName) < 1 {
			c.JSON(http.StatusOK, httpbase.NewBaseRsp(c, httpbase.ParmsError,
				fmt.Sprintf("groupService.StartTunnel, err: %v", "nickName invalid"), nil))
		}
		if err := groupService.StartTunnel(req.GroupName, req.NickName); err != nil {
			c.JSON(http.StatusOK, httpbase.NewBaseRsp(c, httpbase.ParmsError,
				fmt.Errorf("groupService.StartTunnel, err: %v", err).Error(), nil))
			return
		}
	} else if req.Operation == group.OperationStop {
		if err := groupService.StopTunnel(req.GroupName); err != nil {
			c.JSON(http.StatusOK, httpbase.NewBaseRsp(c, httpbase.ParmsError,
				fmt.Errorf("groupService.StopTunnel, err: %v", err).Error(), nil))
			return
		}
	} else {
		c.JSON(http.StatusOK, httpbase.NewBaseRsp(c, httpbase.ParmsError,
			fmt.Errorf("groupService, unsupported operation %v", req.Operation).Error(), nil))
		return
	}

	c.JSON(http.StatusOK, httpbase.NewOkBaseRsp(c, nil))
}
