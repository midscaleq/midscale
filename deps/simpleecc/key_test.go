package simpleecc

import (
	"crypto/elliptic"
	"testing"
)

func TestKeyToHex(t *testing.T) {

	curve := elliptic.P256()
	pri, pub, err := GenKey(curve)
	if err != nil {
		t.Errorf("failed GenKey, err:%v\n", err)
		return
	}

	priHex := PriKeyToHex(pri)
	t.Logf("prikey hex:%v\n", priHex)

	pubHex := PubKeyToHex(curve, pub)
	t.Logf("pubkey hex:%v\n", pubHex)

}

func TestHexToKey(t *testing.T) {

	curve := elliptic.P256()
	pri, pub, err := GenKey(curve)
	if err != nil {
		t.Errorf("failed GenKey, err:%v\n", err)
		return
	}

	priHex := PriKeyToHex(pri)
	// t.Logf("prikey hex:%v\n", priHex)

	pubHex := PubKeyToHex(curve, pub)
	// t.Logf("pubkey hex:%v\n", pubHex)

	pri1, pub1, err := HexToPriKey(curve, priHex)
	if err != nil {
		t.Errorf("failed HexToPriKey, err:%v\n", err)
	}

	// t.Logf("pri:%+v\n", pri)
	// t.Logf("pri1:%+v\n", pri1)

	t.Logf("pub:%+v\n", pub)
	t.Logf("pub1:%+v\n", pub1)

	if !pri.Equal(pri1) {
		t.Errorf("HexToPriKey prikey not Equal\n")
	}

	if !pub.Equal(&pub1) {
		t.Errorf("HexToPriKey pubkey not Equal\n")
	}

	pub2, err := HexToPubKey(curve, pubHex)
	if err != nil {
		t.Errorf("failed HexToPriKey, err:%v\n", err)
	}
	t.Logf("pub2:%+v\n", pub2)
	if !pub.Equal(&pub2) {
		t.Errorf("HexToPubKey pubkey not Equal\n")
	}
}
