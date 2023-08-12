package key

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"fmt"
	"midscale/midscale/app/data/model/define"
	"midscale/midscale/app/data/model/transport/key"
	"midscale/midscale/deps/mskit/file"
	"midscale/midscale/deps/simpleecc"
)

const (
	KeyFileName = "key.data"
)

var priKeyHex string
var pubKeyHex string

var priKey *ecdsa.PrivateKey
var pubKey ecdsa.PublicKey

func Curve() elliptic.Curve {
	return elliptic.P256()
}

func SetUp() error {
	curve := elliptic.P256()

	fileFullName, err := getFileFullName()
	if err != nil {
		return fmt.Errorf("SetUp, getFileFullName err: %v", err)
	}

	if file.IsFileExist(fileFullName) {
		var transkey key.TransportKey
		err = file.ReadFileJsonToObject(fileFullName, &transkey)
		if err != nil {
			return fmt.Errorf("SetUp, ReadFileJsonToObject err: %v", err)
		}
		priKeyHex = transkey.PriKey
		pubKeyHex = transkey.PubKey

		priKey, pubKey, err = simpleecc.HexToPriKey(curve, priKeyHex)
		return err
	}

	pri, pub, err := simpleecc.GenKey(curve)
	if err != nil {
		return fmt.Errorf("SetUp, GenKey err: %v", err)
	}
	transkey := key.TransportKey{PriKey: simpleecc.PriKeyToHex(pri), PubKey: simpleecc.PubKeyToHex(curve, pub)}
	err = file.WriteToFileAsJson(fileFullName, transkey, "  ", false)
	if err != nil {
		return fmt.Errorf("SetUp, WriteToFileAsJson, file:%v, err: %v", fileFullName, err)
	}

	priKey = pri
	pubKey = pub

	priKeyHex = transkey.PriKey
	pubKeyHex = transkey.PubKey

	return nil
}

func getFileFullName() (string, error) {
	p, err := file.GetAppData(define.App)
	if err != nil {
		return "", fmt.Errorf("getFileFullName, GetAppData err: %v", err)
	}
	fileFullName := file.AddPathSepIfNeed(p) + KeyFileName
	return fileFullName, nil
}

func GetTransportPriKey() *ecdsa.PrivateKey {
	return priKey
}

func GetTransportPubKey() ecdsa.PublicKey {
	return pubKey
}

func GetTransportPriKeyHex() string {
	return priKeyHex
}

func GetTransportPubKeyHex() string {
	return pubKeyHex
}
