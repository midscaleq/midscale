package utils

import (
	"fmt"
	"net"
)

// cidrs := []string{"10.1.0.0/16", "10.1.0.2/16", "10.1.0.255/16", "10.1.1.0/16"}
// newCIDR := generateNewCIDR(cidrs)
func GenerateNewCIDR(cidrs []string) (string, error) {
	if len(cidrs) < 1 {
		return "", fmt.Errorf("input array len is 0")
	}

	_, subnet, err := net.ParseCIDR(cidrs[0])
	if err != nil {
		return "", err
	}

	ip := subnet.IP

	for i := 0; i < 255*255*255; i++ {
		incrementIP(&ip)

		ones, _ := subnet.Mask.Size()
		newCIDR := fmt.Sprintf("%s/%v", ip.String(), ones)

		if !isCIDROverlapping(newCIDR, cidrs) {
			return newCIDR, nil
		}
	}

	return "", fmt.Errorf("not found a valid one")
}

func incrementIP(ip *net.IP) {
	for i := len(*ip) - 1; i >= 0; i-- {
		(*ip)[i]++
		if (*ip)[i] == 0 || (*ip)[i] == 255 {
			(*ip)[i] = 1
			continue
		}
		if (*ip)[i] > 0 {
			break
		}
	}
}

func isCIDROverlapping(cidr string, cidrs []string) bool {
	newNet, _, _ := net.ParseCIDR(cidr)
	for _, c := range cidrs {
		net, _, _ := net.ParseCIDR(c)
		if newNet.Equal(net) {
			return true
		}
	}
	return false
}
