package wgconf

import (
	"fmt"
	"testing"
)

func TestParseWireGuardConfig(t *testing.T) {

	var confData = `[Interface]
PrivateKey = %s
ListenPort = 50521
Address = 10.10.10.1/16
DNS = 1.1.1.1, 114.114.114.114, 8.8.8.8

[Peer]
PublicKey = %s
AllowedIPs = 10.10.10.0/24, 192.168.31.0/24, 192.168.1.0/24
PersistentKeepalive = 10
`

	var privateKey = `EIBGhLjHv2H2yXeFS3RK9oCwGxSj8oSitwIgifjA034=`
	var publicKey = `dB2CIP9bxrVHyrA8l/Pp2LtoYteepvPykE1sDdZL1z0=`

	confData = fmt.Sprintf(confData, privateKey, publicKey)

	// Parse the config file into a WireGuardConfig instance
	config, err := ParseWireGuardConfig(string(confData))
	if err != nil {
		t.Errorf("failed to parse config file: %v", err)
	}

	if config.Interface.PrivateKey != privateKey {
		t.Errorf("HexToPubKey pubkey not Equal")
	}
	if len(config.Peers) < 1 {
		t.Errorf("len(config.Peers)<1")
	}
	if config.Peers[0].PublicKey != publicKey {
		t.Errorf("HexToPubKey pubkey not Equal\n")
	}

	config.Peers = append(config.Peers, config.Peers[0])
	if len(config.Peers) < 2 {
		t.Errorf("len(config.Peers){%v}<2", len(config.Peers))
	}
	if config.Peers[1].PublicKey != publicKey {
		t.Errorf("HexToPubKey pubkey not Equal\n")
	}

	t.Logf("WireGuard config: %+v", config)
}
