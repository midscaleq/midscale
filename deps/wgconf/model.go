package wgconf

import (
	"net"
)

type Interface struct {
	Address    *net.IPNet
	PrivateKey string
	ListenPort int
	DNS        string
}

type Peer struct {
	PublicKey           string
	AllowedIPs          []*net.IPNet
	Endpoint            string
	PersistentKeepalive int
}

type WireGuardConfig struct {
	Interface *Interface
	Peers     []*Peer
}
