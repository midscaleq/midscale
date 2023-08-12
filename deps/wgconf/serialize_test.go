package wgconf

import (
	"net"
	"testing"
)

func TestSerializeWireGuardConfig(t *testing.T) {

	var privateKey = `EIBGhLjHv2H2yXeFS3RK9oCwGxSj8oSitwIgifjA034=`
	var publicKey = `dB2CIP9bxrVHyrA8l/Pp2LtoYteepvPykE1sDdZL1z0=`

	// Read the WireGuard config file

	address := &net.IPNet{
		IP:   net.ParseIP("10.10.10.5"),
		Mask: net.CIDRMask(16, 32),
	}

	allowedIPs := []*net.IPNet{{
		IP:   net.ParseIP("10.10.0.0"),
		Mask: net.CIDRMask(16, 32),
	}, {
		IP:   net.ParseIP("192.168.1.0"),
		Mask: net.CIDRMask(16, 32),
	},
	}

	config := WireGuardConfig{
		Interface: &Interface{
			PrivateKey: privateKey,
			Address:    address,
			ListenPort: 50512,
			DNS:        `1.1.1.1, 114.114.114.114, 8.8.8.8`,
		},
		Peers: []*Peer{
			{
				PublicKey:           publicKey,
				AllowedIPs:          allowedIPs,
				Endpoint:            "192.168.1.2:51820",
				PersistentKeepalive: 25,
			},
			{
				PublicKey:  publicKey,
				AllowedIPs: allowedIPs,
				// Endpoint:            "example.com:51820",
				// PersistentKeepalive: 15,
			},
		},
	}

	confData, _ := SerializeWireGuardConfig(&config)
	t.Logf("WireGuard config: \n%+v", config)

	t.Logf("WireGuard confData: \n%+v", string(confData))
}
