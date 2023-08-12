package simpleecc

import (
	"crypto/ecdsa"
	"crypto/rand"

	"github.com/ethereum/go-ethereum/crypto/ecies"
)

func Encrypt(publicKey *ecdsa.PublicKey, message []byte) ([]byte, error) {

	pub := ecies.ImportECDSAPublic(publicKey)
	ciphertext, err := ecies.Encrypt(rand.Reader, pub, message, nil, nil)
	if err != nil {
		return nil, err
	}
	return ciphertext, nil
}

func Decrypt(privateKey *ecdsa.PrivateKey, ciphertext []byte) ([]byte, error) {
	pri := ecies.ImportECDSA(privateKey)

	plaintext, err := pri.Decrypt(ciphertext, nil, nil)
	if err != nil {
		return nil, err
	}
	return plaintext, nil
}
