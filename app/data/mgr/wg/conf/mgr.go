package conf

import (
	"fmt"
	"midscale/midscale/app/data/model/define"
	"midscale/midscale/deps/mskit/file"
	"midscale/midscale/deps/mskit/net/utils"
	"midscale/midscale/deps/wgconf"
	"net"
	"sync"
)

const (
	fileNameExt = ".conf"
)

var lock sync.Mutex

func NewConf(groupName, nickName, localIP string,
	oneMaskLength, listenPort int,
	dns string) error {

	lock.Lock()
	defer lock.Unlock()

	_, ipNet, err := utils.CIDRStringToIpNet(fmt.Sprintf("%v/%v", localIP, oneMaskLength))
	if err != nil {
		return fmt.Errorf("NewConf, CIDRStringToIpNet err: %v", err)
	}

	pri, _, err := wgconf.GenPriKey()
	if err != nil {
		return fmt.Errorf("NewConf, GenPriKey err: %v", err)
	}

	config := &wgconf.WireGuardConfig{
		Interface: &wgconf.Interface{
			Address:    ipNet,
			PrivateKey: pri,
			ListenPort: listenPort,
			DNS:        dns,
		},
	}

	return write(groupName, nickName, config)
}

func GetConfBuf(groupName, nickName string) ([]byte, error) {
	lock.Lock()
	defer lock.Unlock()

	fileFullName, err := ConfFileName(groupName, nickName)
	if err != nil {
		return nil, err
	}
	if !file.IsFileExist(fileFullName) {
		return nil, fmt.Errorf("doGetConf, file %v not exist", fileFullName)
	}
	return file.ReadFromFile(fileFullName)
}

func GetConf(groupName, nickName string) (*wgconf.WireGuardConfig, error) {
	lock.Lock()
	defer lock.Unlock()

	return doGetConf(groupName, nickName)
}

func doGetConf(groupName, nickName string) (*wgconf.WireGuardConfig, error) {

	fileFullName, err := ConfFileName(groupName, nickName)
	if err != nil {
		return nil, err
	}
	if !file.IsFileExist(fileFullName) {
		return nil, fmt.Errorf("doGetConf, file %v not exist", fileFullName)
	}

	b, err := file.ReadFromFile(fileFullName)
	if err != nil {
		return nil, err
	}

	return wgconf.ParseWireGuardConfig(string(b))
}

func ConfFileName(groupName, nickName string) (string, error) {
	fileFullName, err := file.GenFullFileNameByAppData(define.App, groupName+fileNameExt, groupName, nickName)
	if err != nil {
		return "", fmt.Errorf("ConfFileName, GenFullFileNameByAppData err: %v", err)
	}
	return fileFullName, nil
}

func ConfExist(groupName, nickName string) (bool, error) {
	lock.Lock()
	defer lock.Unlock()

	fileFullName, err := ConfFileName(groupName, nickName)
	if err != nil {
		return false, err
	}
	return file.IsFileExist(fileFullName), nil
}

func UpdateLocalIp(groupName, nickName string, localIP string) error {
	lock.Lock()
	defer lock.Unlock()

	config, err := doGetConf(groupName, nickName)
	if err != nil {
		return err
	}

	_, ipNet, err := utils.CIDRStringToIpNet(localIP)
	if err != nil {
		return fmt.Errorf("UpdateLocalIp, CIDRStringToIpNet err: %v", err)
	}

	config.Interface.Address = ipNet
	return write(groupName, nickName, config)
}

func AppendPeerStr(groupName, nickName, endpoint, publicKey string, allowedIPsStr []string, persistentKeepalive int) error {

	allowedIPs := make([]*net.IPNet, 0)
	for _, v := range allowedIPsStr {
		_, subnet, err := utils.CIDRStringToIpNet(v)
		if err != nil {
			return fmt.Errorf("AppendPeerStr, CIDRStringToIpNet err: %v", err)
		}
		allowedIPs = append(allowedIPs, subnet)
	}

	return AppendPeer(groupName, nickName, endpoint, publicKey, allowedIPs, persistentKeepalive)
}

func AppendPeer(groupName, nickName, endpoint, publicKey string, allowedIPs []*net.IPNet, persistentKeepalive int) error {
	lock.Lock()
	defer lock.Unlock()

	config, err := doGetConf(groupName, nickName)
	if err != nil {
		return err
	}
	for _, p := range config.Peers {
		if p.PublicKey == publicKey {
			return fmt.Errorf("publicKey already exist")
		}
	}

	peer := &wgconf.Peer{PublicKey: publicKey,
		Endpoint:            endpoint,
		AllowedIPs:          allowedIPs,
		PersistentKeepalive: persistentKeepalive}
	config.Peers = append(config.Peers, peer)
	return write(groupName, nickName, config)
}

func GetPeers(groupName, nickName, serverEndpoint string, serverPersistentKeepalive int) ([]*wgconf.Peer, error) {
	lock.Lock()
	defer lock.Unlock()

	config, err := doGetConf(groupName, nickName)
	if err != nil {
		return nil, err
	}

	configIncludeServer := make([]*wgconf.Peer, 0)
	serverPubKey, err := wgconf.GetPubKey(config.Interface.PrivateKey)
	if err != nil {
		return nil, err
	}
	replacedServerLocalIPString, err := utils.ReplaceCIDRStringWithNewNetMask(config.Interface.Address.String(), "32")
	if err != nil {
		return nil, err
	}
	_, replacedServerLocalIPNet, err := utils.CIDRStringToIpNet(replacedServerLocalIPString)
	if err != nil {
		return nil, err
	}
	allowedIPs := []*net.IPNet{replacedServerLocalIPNet}
	serverPeer := &wgconf.Peer{PublicKey: serverPubKey, AllowedIPs: allowedIPs,
		Endpoint: serverEndpoint, PersistentKeepalive: serverPersistentKeepalive}

	configIncludeServer = append(configIncludeServer, serverPeer)
	configIncludeServer = append(configIncludeServer, config.Peers...)
	return configIncludeServer, nil
}

func GetPubKey(groupName, nickName string) (string, error) {
	lock.Lock()
	defer lock.Unlock()

	config, err := doGetConf(groupName, nickName)
	if err != nil {
		return "", err
	}
	pri := config.Interface.PrivateKey
	if len(pri) < 1 {
		return "", fmt.Errorf("pri is empty")
	}
	return wgconf.GetPubKey(pri)
}

func write(groupName, nickName string, config *wgconf.WireGuardConfig) error {

	b, err := wgconf.SerializeWireGuardConfig(config)
	if err != nil {
		return fmt.Errorf("write, SerializeWireGuardConfig err: %v", err)
	}

	fileFullName, err := ConfFileName(groupName, nickName)
	if err != nil {
		return err
	}

	if err := file.WriteToFile(fileFullName, b, true); err != nil {
		return fmt.Errorf("failed write, file:%v, err:%v", fileFullName, err)
	}
	return nil
}
