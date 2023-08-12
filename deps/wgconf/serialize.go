package wgconf

import (
	"fmt"
	"strings"
)

func SerializeWireGuardConfig(config *WireGuardConfig) ([]byte, error) {

	content := "[Interface]\n"
	content += fmt.Sprintf("PrivateKey = %s\n", config.Interface.PrivateKey)
	if len(config.Interface.Address.String()) > 0 {
		content += fmt.Sprintf("Address = %s\n", config.Interface.Address.String())
	}
	if config.Interface.ListenPort > 0 {
		content += fmt.Sprintf("ListenPort = %d\n", config.Interface.ListenPort)
	}
	if len(config.Interface.DNS) > 0 {
		content += fmt.Sprintf("DNS = %s\n", config.Interface.DNS)
	}
	content += "\n"

	for _, peer := range config.Peers {
		allowedIPs := make([]string, 0)
		for _, ipNet := range peer.AllowedIPs {
			one, _ := ipNet.Mask.Size()
			s := fmt.Sprintf("%s/%d", ipNet.IP, one)
			allowedIPs = append(allowedIPs, s)
		}
		allowedIPsJoined := strings.Join(allowedIPs, ",")

		content += "[Peer]\n"
		content += fmt.Sprintf("PublicKey = %s\n", peer.PublicKey)
		if len(allowedIPsJoined) > 0 {
			content += fmt.Sprintf("AllowedIPs = %s\n", allowedIPsJoined)
		}
		if len(peer.Endpoint) > 0 {
			content += fmt.Sprintf("Endpoint = %s\n", peer.Endpoint)
		}
		if peer.PersistentKeepalive > 0 {
			content += fmt.Sprintf("PersistentKeepalive = %d\n\n", peer.PersistentKeepalive)
		}
	}

	return []byte(content), nil
}
