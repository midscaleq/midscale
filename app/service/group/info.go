package group

import (
	"fmt"
	groupinfo "midscale/midscale/app/data/mgr/group/info"
	"midscale/midscale/app/data/model/group/info"
)

type InfoService struct {
}

func NewInfoService() *InfoService {
	return &InfoService{}
}

func (infoService *InfoService) GetInfo(groupName, nickName string) (*info.GroupInfo, error) {

	info, err := groupinfo.GetGroupInfo(groupName, nickName)
	if err != nil {
		return nil, fmt.Errorf("InfoService.GetInfo, GetGroupInfo err: %v", err)
	}

	return info, nil
}
