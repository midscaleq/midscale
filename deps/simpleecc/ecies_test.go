package simpleecc

import (
	"crypto/elliptic"
	"testing"
)

func TestEncrypt(t *testing.T) {

	curve := elliptic.P256()
	pri, pub, err := GenKey(curve)
	if err != nil {
		t.Errorf("failed GenKey, err:%v\n", err)
		return
	}

	message := "hello, test ecies!"
	ciphertext, err := Encrypt(&pub, []byte(message))
	if err != nil {
		t.Errorf("failed Encrypt, err:%v\n", err)
		return
	}

	decryptedText, err := Decrypt(pri, ciphertext)
	if err != nil {
		t.Errorf("failed Decrypt, err:%v\n", err)
		return
	}

	if string(message) == string(decryptedText) {
		// t.Logf("加解密通过, ciphertext：%v\n", ciphertext)
	} else {
		t.Errorf("message %v != decryptedText %v ", string(message), string(decryptedText))
	}
}
