package wgconf

import "golang.zx2c4.com/wireguard/wgctrl/wgtypes"

func GenPriKey() (string, string, error) {
	// 生成私钥
	privateKey, err := wgtypes.GeneratePrivateKey()
	if err != nil {
		return "", "", err
	}
	return privateKey.String(), privateKey.PublicKey().String(), nil
}

func GetPubKey(priKey string) (string, error) {
	key, err := wgtypes.ParseKey(priKey)
	if err != nil {
		return "", err
	}
	return key.PublicKey().String(), nil
}

func GetKey(priKey string) (string, string, error) {
	key, err := wgtypes.ParseKey(priKey)
	if err != nil {
		return "", "", err
	}
	return key.String(), key.PublicKey().String(), nil
}
