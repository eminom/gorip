package xau

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
)

var (
	hmac []byte
)

func init() {
	buf := make([]byte, 1024)
	nl, ne := hex.Decode(buf, []byte(_hashSecret))
	if nil != ne {
		panic(ne)
	}
	hmac = buf[:nl]
}

func CheckUserPasswd(u, p string) bool {
	hd := sha256.Sum256([]byte(u + p))
	return 0 == bytes.Compare(hmac, hd[:])
}
