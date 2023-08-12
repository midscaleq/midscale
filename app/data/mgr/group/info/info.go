package info

import (
	"fmt"
	"midscale/midscale/app/data/model/define"
	"midscale/midscale/app/data/model/group/info"
	"midscale/midscale/deps/mskit/file"
	"midscale/midscale/deps/mskit/net/utils"
	"sync"
)

const (
	fileNameExt = ".group"
)

var lock sync.Mutex

func init() {
}

func NewGroupInfo(groupInfo *info.GroupInfo) error {
	lock.Lock()
	defer lock.Unlock()

	if err := write(groupInfo); err != nil {
		return err
	}
	return nil
}

func VerifyJoinPassword(groupName, nickName, pwd string) (bool, error) {
	groupInfo, err := read(groupName, nickName)
	if err != nil {
		return false, err
	}

	if groupInfo == nil {
		return false, nil
	}

	return groupInfo.ServerJoinPassword == pwd, nil
}

func AppendJoiner(groupName, nickName string, joiner *info.JoinerInfo) error {
	lock.Lock()
	defer lock.Unlock()

	groupInfo, err := read(groupName, nickName)
	if err != nil {
		return err
	}
	if groupInfo == nil {
		return fmt.Errorf("group file is empty")
	}

	if groupInfo.JoinerInfos == nil {
		groupInfo.JoinerInfos = make([]*info.JoinerInfo, 0)
	}

	joinerLocalIPs := make([]string, 0)
	joinerLocalIPs = append(joinerLocalIPs, groupInfo.ServerLocalIPCIDR)
	for _, cli := range groupInfo.JoinerInfos {
		if cli.JoinerPubKey == joiner.JoinerPubKey {
			return fmt.Errorf("AppendJoiner, joiner exist! JoinerPubKey:%v", cli.JoinerPubKey)
		}
		joinerLocalIPs = append(joinerLocalIPs, cli.JoinerLocalIP)
	}

	// auto assign new ip
	if len(joiner.JoinerLocalIP) < 1 {
		joiner.JoinerLocalIP, err = utils.GenerateNewCIDR(joinerLocalIPs)
		if err != nil {
			return fmt.Errorf("AppendJoiner, GenerateNewCIDR err: %v", err)
		}
	}

	groupInfo.JoinerInfos = append(groupInfo.JoinerInfos, joiner)

	if err := write(groupInfo); err != nil {
		return err
	}
	return nil
}

func GroupExist(groupName string) (bool, error) {
	lock.Lock()
	defer lock.Unlock()

	p, err := file.GenPathByAppData(define.App, groupName)
	if err != nil {
		return false, fmt.Errorf("GroupExist, GenPathByAppData err: %v", err)
	}

	if file.IsFileExist(p) {
		return true, nil
	}
	return false, nil
}

func GroupJoined(groupName, nickName string) bool {
	_, err := GetGroupInfo(groupName, nickName)
	return err == nil
}

func GetGroupInfo(groupName, nickName string) (*info.GroupInfo, error) {
	lock.Lock()
	defer lock.Unlock()

	return read(groupName, nickName)
}

func write(g *info.GroupInfo) error {
	fileFullName, err := file.GenFullFileNameByAppData(define.App,
		g.GroupName+fileNameExt, g.GroupName, g.ServerNickName)
	if err != nil {
		return fmt.Errorf("write, GenFullFileNameByAppData err: %v", err)
	}

	if err := file.WriteToFileAsJson(fileFullName, g, "  ", true); err != nil {
		return fmt.Errorf("failed WriteToFileAsJson, file:%v, err:%v", fileFullName, err)
	}
	return nil
}

func read(groupName, nickName string) (*info.GroupInfo, error) {
	fileFullName, err := file.GenFullFileNameByAppData(define.App,
		groupName+fileNameExt, groupName, nickName)
	if err != nil {
		return nil, fmt.Errorf("read, GenFullFileNameByAppData err: %v", err)
	}

	var g info.GroupInfo
	if err = file.ReadFileJsonToObject(fileFullName, &g); err != nil {
		return nil, fmt.Errorf("failed ReadFileJsonToObject, file:%v, err:%v", fileFullName, err)
	}
	return &g, nil
}
