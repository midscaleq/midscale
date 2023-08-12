package simpleecc

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"math/big"
)

func GenKey(curve elliptic.Curve) (*ecdsa.PrivateKey, ecdsa.PublicKey, error) {
	var publicKey ecdsa.PublicKey
	privateKey, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err == nil {
		publicKey = privateKey.PublicKey
	}
	return privateKey, publicKey, err
}

func PriKeyToBase64(privateKey *ecdsa.PrivateKey) string {
	privateKeyBytes := privateKey.D.Bytes()
	// fmt.Printf("outputPriKey : %x\n", privateKey.D.Bytes())
	privateKeyBase64 := base64.StdEncoding.EncodeToString(privateKeyBytes)
	// fmt.Println("outputPriKey privateKeyBase64:", privateKeyBase64)
	return privateKeyBase64
}

func PriKeyToHex(privateKey *ecdsa.PrivateKey) string {
	privateKeyBytes := privateKey.D.Bytes()
	// fmt.Printf("outputPriKey : %x\n", privateKey.D.Bytes())
	privateKeyHex := hex.EncodeToString(privateKeyBytes)
	// fmt.Println("outputPriKey  privateKeyBase64:", privateKeyBase64)
	return privateKeyHex
}

func HexToPriKey(curve elliptic.Curve, privateKeyHex string) (*ecdsa.PrivateKey, ecdsa.PublicKey, error) {
	var publicKey ecdsa.PublicKey
	privateKeyBytes, err := hex.DecodeString(privateKeyHex)
	if err != nil {
		return nil, publicKey, err
	}

	// secp256r1
	privateKey := &ecdsa.PrivateKey{D: new(big.Int).SetBytes(privateKeyBytes), PublicKey: ecdsa.PublicKey{Curve: curve}}

	// fmt.Printf("ConstructPriKey : %x\n", privateKey.D.Bytes())

	// gen pub key
	var x, y *big.Int
	x, y = curve.ScalarBaseMult(privateKeyBytes)

	// 取出公钥点的 x 和 y 坐标信息
	// x := privateKey.PublicKey.X.Bytes()
	// y := privateKey.PublicKey.Y.Bytes()

	// 使用 ecdsa.PublicKey 结构体构造完整的公钥
	publicKey = ecdsa.PublicKey{Curve: curve, X: x, Y: y}
	// 使用 x、y 和椭圆曲线参数重新构造完整的公钥对象
	// publicKey := &ecdsa.PublicKey{
	// 	Curve: curve,
	// 	X:     new(big.Int).SetBytes(x),
	// 	Y:     new(big.Int).SetBytes(y),
	// }

	privateKey.PublicKey = publicKey

	return privateKey, publicKey, nil
}

func PubKeyToBase64(curve elliptic.Curve, publicKey ecdsa.PublicKey) string {
	// 将公钥转换成字节数组
	// publicKeyBytes := elliptic.Marshal(curve, publicKey.X, publicKey.Y)
	publicKeyBytes := elliptic.MarshalCompressed(curve, publicKey.X, publicKey.Y)
	// 对字节数组进行 Base64 编码
	publicKeyBase64 := base64.StdEncoding.EncodeToString(publicKeyBytes)
	// fmt.Println("公钥：", publicKeyBase64)
	return publicKeyBase64
}

func PubKeyToHex(curve elliptic.Curve, publicKey ecdsa.PublicKey) string {
	// 将公钥转换成字节数组
	// publicKeyBytes := elliptic.Marshal(curve, publicKey.X, publicKey.Y)
	publicKeyBytes := elliptic.MarshalCompressed(curve, publicKey.X, publicKey.Y)
	publicKeyHex := hex.EncodeToString(publicKeyBytes)
	// fmt.Println("公钥：", publicKeyBase64)
	return publicKeyHex
}

func HexToPubKey(curve elliptic.Curve, publicKeyHex string) (ecdsa.PublicKey, error) {
	var publicKey ecdsa.PublicKey
	// 解码公钥
	publicKeyBytes, err := hex.DecodeString(publicKeyHex)
	if err != nil {
		return publicKey, err
	}

	// 使用曲线构造公钥
	// x, y := elliptic.Unmarshal(curve, publicKeyBytes)
	x, y := elliptic.UnmarshalCompressed(curve, publicKeyBytes)
	publicKey = ecdsa.PublicKey{Curve: curve, X: x, Y: y}

	// 打印公钥
	// fmt.Printf("x: %x\n", publicKey.X)
	// fmt.Printf("y: %x\n", publicKey.Y)
	return publicKey, nil
}

func PriKeyToPEM(priKey *ecdsa.PrivateKey) ([]byte, error) {
	// x509
	derText, err := x509.MarshalECPrivateKey(priKey)
	if err != nil {
		return nil, err
	}
	// pem block
	block := &pem.Block{
		Type:  "ecdsa private key",
		Bytes: derText,
	}

	var buf bytes.Buffer
	err = pem.Encode(&buf, block)
	return buf.Bytes(), err
}

func PubKeyToPEM(pubKey ecdsa.PublicKey) ([]byte, error) {
	derText, err := x509.MarshalPKIXPublicKey(&pubKey)
	if err != nil {
		return nil, err
	}
	block := &pem.Block{
		Type:  "ecdsa public key",
		Bytes: derText,
	}
	var buf bytes.Buffer
	err = pem.Encode(&buf, block)
	return buf.Bytes(), err
}
