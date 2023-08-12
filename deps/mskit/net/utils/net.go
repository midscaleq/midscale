package utils

import (
	"fmt"
	"log"
	"net"
	"strings"
)

// IpStringToIpNet as a CIDR notation IP address and prefix length,
// like "192.0.2.0/24" or "2001:db8::/32"
//
// It returns the IP address and the network implied by the IP and
// prefix length.
// For example, ParseCIDR("192.0.2.1/24") returns the IP address
// 192.0.2.1 and the network 192.0.2.0/24.
func CIDRStringToIpNet(s string) (net.IP, *net.IPNet, error) {
	ip, subnet, err := net.ParseCIDR(s)
	if err != nil {
		log.Printf("CIDRStringToIpNet, s:%+v, er:%v\n", s, err)
		return ip, nil, err
	}

	return ip, &net.IPNet{IP: ip, Mask: subnet.Mask}, nil
}
func ReplaceCIDRStringWithNewNetMask(cidr, newMask string) (string, error) {
	// log.Printf("ReplaceCIDRStringWithNewNetMask, cidr:%v, newMask:%v", cidr, newMask)
	i := strings.Index(cidr, "/")
	if i < 1 {
		return "", fmt.Errorf("cidr format error, cidr:%v", cidr)
	}

	return cidr[:i] + "/" + newMask, nil
}
