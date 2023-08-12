package simpleecc

import (
	"crypto/elliptic"
	"testing"
)

func TestECDH(t *testing.T) {

	curve := elliptic.P256()
	pri, pub, err := GenKey(curve)
	if err != nil {
		t.Errorf("failed GenKey, err:%v\n", err)
		return
	}

	pri2, pub2, err := GenKey(curve)
	if err != nil {
		t.Errorf("failed GenKey, err:%v\n", err)
		return
	}

	sharedKey1, err := ECDH(curve, pri, &pub2)
	if err != nil {
		t.Errorf("failed ECDH, err:%v\n", err)
		return
	}

	sharedKey2, err := ECDH(curve, pri2, &pub)
	if err != nil {
		t.Errorf("failed ECDH, err:%v\n", err)
		return
	}

	if string(sharedKey1) == string(sharedKey2) {
		// t.Logf("共享密钥为：%s\n", sharedKey1)
	} else {
		t.Errorf("sharedKey1 %v != sharedKey2 %v ", string(sharedKey1), string(sharedKey2))
	}
}
