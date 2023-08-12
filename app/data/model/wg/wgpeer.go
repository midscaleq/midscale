package wg

type WgPeerComment struct {
	NickName string `json:"nickName,omitempty"`
	PeerIp   string `json:"peerIp"`
	JoinTime string `json:"joinTime"`
}
