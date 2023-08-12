package wgconf

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

func ParseWireGuardConfig(configStr string) (*WireGuardConfig, error) {
	configStr = strings.ReplaceAll(configStr, "\r\n", "\n")
	configStr = strings.ReplaceAll(configStr, "\r", "\n")

	config := &WireGuardConfig{}

	// Split the config file into blocks
	blocks := split(configStr, "\n[")
	// for i, v := range blocks {
	// 	fmt.Printf("block[%v]: {%v}\n", i, v)
	// }

	// Parse the [Interface] block
	intfStart := strings.Index(blocks[0], "\n") + 1
	intfEnd := len(blocks[0]) - 1
	intfBlock := blocks[0][intfStart:intfEnd]
	// fmt.Printf("==== parseInterfaceBlock intfBlock:%+v\n", intfBlock)
	intf, err := parseInterfaceBlock(intfBlock)
	if err != nil {
		// fmt.Printf("parseInterfaceBlock err:%v\n", err)
		return nil, err
	}

	config.Interface = intf

	// Parse the [Peer] blocks
	for _, blockStr := range blocks[1:] {
		blockEnd := len(blockStr) - 2
		blockStr = blockStr[:blockEnd]

		peer, err := parsePeerBlock(blockStr)
		if err != nil {
			return nil, err
		}

		config.Peers = append(config.Peers, peer)
	}

	return config, nil
}

func parseInterfaceBlock(blockStr string) (*Interface, error) {
	blockLines := strings.Split(blockStr, "\n")

	intf := &Interface{}

	for _, line := range blockLines {
		line = strings.ReplaceAll(line, " ", "")
		line = strings.ReplaceAll(line, "\r\n", "\n")
		line = strings.ReplaceAll(line, "\r", "\n")
		line = strings.ReplaceAll(line, "\n", "")
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// fmt.Printf("==== parseInterfaceBlock Trim line:%+v\n", line)

		keyValue := strings.SplitN(line, "=", 2)
		if len(keyValue) != 2 {
			return nil, fmt.Errorf("parseInterfaceBlock invalid key value pair: %s", line)
		}

		switch keyValue[0] {
		case "Address":
			// Parse the Address field as a net.IPNet
			ip, subnet, err := net.ParseCIDR(keyValue[1])
			if err != nil {
				return nil, fmt.Errorf("parseInterfaceBlock invalid Address value: %s", keyValue[1])
			}

			// Set the Address field of the Interface instance
			intf.Address = &net.IPNet{IP: ip, Mask: subnet.Mask}

		case "PrivateKey":
			if len(keyValue[1]) != 44 {
				return nil, fmt.Errorf("invalid PrivateKey value: %s", keyValue[1])
			}

			// Set the PrivateKey field of the Interface instance
			intf.PrivateKey = keyValue[1]

		case "ListenPort":
			// Parse the ListenPort field as an integer
			port, err := strconv.Atoi(keyValue[1])
			if err != nil {
				return nil, fmt.Errorf("parseInterfaceBlock invalid ListenPort value: [%s], err:%v", keyValue[1], err)
			}

			// Set the ListenPort field of the Interface instance
			intf.ListenPort = port

		case "DNS":
			// Set the PrivateKey field of the Interface instance
			intf.DNS = keyValue[1]

		default:
			return nil, fmt.Errorf("unsupported Interface parameter: %s", keyValue[0])
		}
	}

	return intf, nil
}

func parsePeerBlock(blockStr string) (*Peer, error) {
	blockLines := strings.Split(blockStr, "\n")

	peer := &Peer{}

	for _, line := range blockLines {
		line = strings.ReplaceAll(line, " ", "")
		line = strings.ReplaceAll(line, "\r\n", "\n")
		line = strings.ReplaceAll(line, "\r", "\n")
		line = strings.ReplaceAll(line, "\n", "")
		line = strings.TrimSpace(line)
		if line == "" || strings.Contains(line, "Peer]") {
			continue
		}

		// fmt.Printf("==== parsePeerBlock line:{%+v}\n", line)

		keyValue := strings.SplitN(line, "=", 2)
		if len(keyValue) != 2 {
			return nil, fmt.Errorf("invalid key value pair: %s", line)
		}

		switch keyValue[0] {
		case "PublicKey":
			if len(keyValue[1]) < 44 || len(keyValue[1]) > 44 {
				return nil, fmt.Errorf("invalid PublicKey value: %s", keyValue[1])
			}

			// Set the PublicKey field of the Peer instance
			peer.PublicKey = keyValue[1]

		case "AllowedIPs":
			// Parse the AllowedIPs field as a list of IP networks
			allowedIPs, err := parseIPNetList(keyValue[1])
			if err != nil {
				return nil, fmt.Errorf("invalid AllowedIPs value: %s", keyValue[1])
			}

			// Set the AllowedIPs field of the Peer instance
			peer.AllowedIPs = allowedIPs

		case "Endpoint":
			// Set the Endpoint field of the Peer instance
			peer.Endpoint = keyValue[1]

		case "PersistentKeepalive":
			// Parse the PersistentKeepalive field as an integer
			pka, err := strconv.Atoi(keyValue[1])
			if err != nil {
				return nil, fmt.Errorf("invalid PersistentKeepalive value: %s", keyValue[1])
			}

			// Set the PersistentKeepAlive field of the Peer instance
			peer.PersistentKeepalive = pka

		default:
			return nil, fmt.Errorf("unsupported Peer parameter: %s", keyValue[0])
		}
	}

	return peer, nil
}

func parseIPNetList(str string) ([]*net.IPNet, error) {
	// Split the string into IP network strings
	ipStrs := split(str, ",")

	// Parse each IP network string as an net.IPNet
	nets := make([]*net.IPNet, 0, len(ipStrs))

	for _, ipStr := range ipStrs {
		ip, subnet, err := net.ParseCIDR(ipStr)
		if err != nil {
			return nil, err
		}

		nets = append(nets, &net.IPNet{IP: ip, Mask: subnet.Mask})
	}

	return nets, nil
}

func split(str string, sep string) []string {
	// Split the string into fields
	fields := make([]string, 0)
	start := 0
	for i := 0; i < len(str); i++ {
		if i+len(sep) > len(str) {
			break
		}

		if str[i:i+len(sep)] == sep {
			fields = append(fields, str[start:i])
			start = i + len(sep)
		}
	}

	if start < len(str) {
		fields = append(fields, str[start:])
	}

	return fields
}
