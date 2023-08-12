package define

var DefaultListenAddrForRemote string
var DefaultListenAddrForLocal string
var DefaultListenPort int
var DefaultLocalIP string
var DefaultOneMaskLength int
var DefaultDNS string

var ServerHost string
var EndpointHost string
var PersistentKeepalive int

var AutoStartTunnel bool

func init() {
	ServerHost = "http://192.168.99.5"
	EndpointHost = "192.168.99.5"
	PersistentKeepalive = 30

	DefaultListenAddrForRemote = ":17979"
	DefaultListenAddrForLocal = "127.0.0.1:18888"
	DefaultListenPort = 50501
	DefaultLocalIP = "10.0.0.1"
	DefaultOneMaskLength = 16
	DefaultDNS = `1.1.1.1, 114.114.114.114, 8.8.8.8`

	AutoStartTunnel = true
}
