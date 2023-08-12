package json

import (
	"crypto/elliptic"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"midscale/midscale/deps/simpleecc"
)

const (
	hasLen = 4
)

func defaultCurve() elliptic.Curve {
	return elliptic.P256()
}

func DecryptB64ToJson(priKeyHex, base64String, plainHash string, returnedJsonObj interface{}) error {
	// log.Printf("#### DecryptB64ToJson, priKeyHex:%v, base64String:%+v",
	// 	priKeyHex, base64String)

	data, err := base64.StdEncoding.DecodeString(base64String)
	if err != nil {
		return fmt.Errorf("decodeString err: %v", err)
	}

	pri, _, err := simpleecc.HexToPriKey(defaultCurve(), priKeyHex)
	if err != nil {
		return fmt.Errorf("hexToPriKey err: %v", err)
	}

	plain, err := simpleecc.Decrypt(pri, data)
	if err != nil {
		return fmt.Errorf("decrypt err: %v", err)
	}

	plainHashRaw := sha256.Sum256(plain)
	if plainHash != base64.StdEncoding.EncodeToString(plainHashRaw[0:4]) {
		return fmt.Errorf("hash invalid")
	}

	err = json.Unmarshal(plain, returnedJsonObj)
	if err != nil {
		return fmt.Errorf("unmarshal err: %v", err)
	}
	// log.Printf("#### DecryptB64ToJson, priKeyHex:%v, base64String:%+v, returnedJsonObj:%+v",
	// 	priKeyHex, base64String, returnedJsonObj)

	return nil
}

func EncryptJsonToB64(pubKeyHex string, jsonObj interface{}) (string, string, error) {
	// log.Printf("#### EncryptJsonToB64, pubKeyHex:%v, jsonObj:%+v", pubKeyHex, jsonObj)
	pub, err := simpleecc.HexToPubKey(defaultCurve(), pubKeyHex)
	if err != nil {
		return "", "", fmt.Errorf("encryptJsonToB64 err: %v", err)
	}

	data, err := json.Marshal(jsonObj)
	if err != nil {
		return "", "", fmt.Errorf("marshal err: %v", err)
	}

	ciphertext, err := simpleecc.Encrypt(&pub, data)
	if err != nil {
		return "", "", fmt.Errorf("encrypt err: %v", err)
	}

	b64 := base64.StdEncoding.EncodeToString(ciphertext)
	// log.Printf("#### EncryptJsonToB64, pubKeyHex:%v, jsonObj:%+v, b64:%+v",
	// 	pubKeyHex, jsonObj, b64)

	plainHashRaw := sha256.Sum256(data)

	return b64, base64.StdEncoding.EncodeToString(plainHashRaw[0:4]), nil
}
