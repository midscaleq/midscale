package random

import (
	"crypto/rand"
	"encoding/hex"

	uuid "github.com/satori/go.uuid"
)

func Random(byteLen uint16) []byte {

	b := make([]byte, byteLen)
	n, err := rand.Read(b) //在byte切片中随机写入元素
	if err != nil {
		return nil
	}
	if n != int(byteLen) {
		return nil
	}
	return b
}

func RandomHex(byteLen uint16) string {

	b := make([]byte, byteLen)
	n, err := rand.Read(b) //在byte切片中随机写入元素
	if err != nil {
		return ""
	}
	if n != int(byteLen) {
		return ""
	}
	return hex.EncodeToString(b)
}

func GenUUID() string {
	u1 := uuid.NewV4()
	return u1.String()
}
