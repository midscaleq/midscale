package simpleecc

import (
	"crypto/ecdsa"
	"crypto/rand"
	"hash"
)

func SignASN1(priKey *ecdsa.PrivateKey, myHash hash.Hash, message []byte) ([]byte, error) {
	_, err := myHash.Write(message)
	if err != nil {
		return nil, err
	}
	messageHash := myHash.Sum(nil)
	sign, err := ecdsa.SignASN1(rand.Reader, priKey, messageHash)
	return sign, err
}

func VerifyASN1(publicKey *ecdsa.PublicKey, myHash hash.Hash, message []byte, sig []byte) (bool, error) {
	_, err := myHash.Write(message)
	if err != nil {
		return false, err
	}
	messageHash := myHash.Sum(nil)

	verified := ecdsa.VerifyASN1(publicKey, messageHash, sig)
	return verified, nil
}
