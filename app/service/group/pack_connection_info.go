package group

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"midscale/midscale/app/data/mgr/transport/key"
	"midscale/midscale/app/data/model/define"
	"midscale/midscale/deps/mskit/encoding/bytes"
	"midscale/midscale/deps/simpleecc"

	"github.com/mr-tron/base58"
)

func PackConnectionInfo(connInfo interface{}) (string, error) {

	data, err := json.Marshal(connInfo)
	if err != nil {
		return "", fmt.Errorf("marshal err: %v", err)
	}
	sign, err := simpleecc.SignASN1(key.GetTransportPriKey(), sha256.New(), data)
	if err != nil {
		return "", fmt.Errorf("signASN1 err: %v", err)
	}
	version := bytes.Int16ToBytes(define.ConnectionInfoVersion)
	dataLen := bytes.Int16ToBytes(int16(len(data)))
	connectionInfo := append(version, dataLen...)
	connectionInfo = append(connectionInfo, data...)
	connectionInfo = append(connectionInfo, sign...)
	fmt.Printf("\n\n PackConnectionInfo, data:%v, sign:%v\n\n", data, sign)
	return base58.Encode(connectionInfo), nil
}

func UnpackConnectionInfo(connectionInfo string) ([]byte, []byte, error) {

	decodedReq, err := base58.Decode(connectionInfo)
	if err != nil {
		return nil, nil, fmt.Errorf("failed base58.Decode %v", err)
	}

	if len(decodedReq) < define.ConnectionInfoVersionLen+define.ConnectionInfoSignLen {
		return nil, nil, fmt.Errorf("len(decodedReq)=%v error", len(decodedReq))
	}

	offset := 0
	version := decodedReq[offset : offset+define.ConnectionInfoVersionLen]
	v := bytes.BytesToInt16(version)
	if v != define.ConnectionInfoVersion {
		return nil, nil, fmt.Errorf("ConnectionInfoVersion=%v error", v)
	}
	offset += define.ConnectionInfoVersionLen
	dataLen := decodedReq[offset : offset+define.ConnectionInfoDataLenLen]
	l := bytes.BytesToInt16(dataLen)
	offset += len(dataLen)
	data := decodedReq[offset : offset+int(l)]
	offset += int(l)
	sign := decodedReq[offset:]

	return data, sign, nil
}

func VerifySign(pubKeyHex string, data, sign []byte) (bool, error) {

	publicKey, err := simpleecc.HexToPubKey(key.Curve(), pubKeyHex)
	if err != nil {
		return false, fmt.Errorf("failed HexToPubKey %v", err)
	}
	return simpleecc.VerifyASN1(&publicKey, sha256.New(), data, sign)
}
