package simpleecc

import (
	"crypto/elliptic"
	"crypto/sha256"
	"testing"
)

func TestSign(t *testing.T) {

	curve := elliptic.P256()
	pri, pub, err := GenKey(curve)
	if err != nil {
		t.Errorf("failed GenKey, err:%v\n", err)
		return
	}

	message := "message to sign!"
	sig, err := SignASN1(pri, sha256.New(), []byte(message))
	if err != nil {
		t.Errorf("failed Sign, err:%v\n", err)
		return
	}

	verified, err := VerifyASN1(&pub, sha256.New(), []byte(message), sig)
	if err != nil {
		t.Errorf("failed Verify, err:%v\n", err)
		return
	}

	if !verified {
		t.Errorf("Verify not pass!")
	}
}
