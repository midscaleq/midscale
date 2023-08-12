package simpleecc

import (
	"crypto/ecdh"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"fmt"
)

func ECDH(curve elliptic.Curve, privKey *ecdsa.PrivateKey, remote *ecdsa.PublicKey) ([]byte, error) {
	// 计算共享密钥
	x, _ := curve.ScalarMult(remote.X, remote.Y, privKey.D.Bytes())

	return x.MarshalText()
}

func ECDH2(k *ecdh.PrivateKey, remote *ecdh.PublicKey) {
	// 生成椭圆曲线
	curve := elliptic.P256()

	// 生成私钥
	privKeyA, _ := ecdsa.GenerateKey(curve, rand.Reader)

	// 生成公钥
	pubKeyA := privKeyA.PublicKey

	// 生成私钥
	privKeyB, _ := ecdsa.GenerateKey(curve, rand.Reader)

	// 生成公钥
	pubKeyB := privKeyB.PublicKey

	// 计算共享密钥
	x, _ := curve.ScalarMult(pubKeyA.X, pubKeyA.Y, privKeyB.D.Bytes())
	y, _ := curve.ScalarMult(pubKeyB.X, pubKeyB.Y, privKeyA.D.Bytes())

	if x.Cmp(y) == 0 {
		fmt.Println("ECDH key exchange success!")
	} else {
		fmt.Println("ECDH key exchange failed!")
	}
}
