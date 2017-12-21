package blocker

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"io"
	"strings"
	"unsafe"
)

//
func decryptEncoded(hexStr, keyStr string) ([]byte, error) {
	hexStr = strings.TrimSpace(hexStr)
	encbuffer, err := hex.DecodeString(hexStr)
	if err != nil {
		return nil, err
	}
	iv := encbuffer[:aes.BlockSize]
	out := encbuffer[aes.BlockSize:]

	realKey := doHash(keyStr)
	block, ec := aes.NewCipher(padKey(realKey))
	if ec != nil {
		return nil, ec
	}
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(out, out)
	usize := binary.BigEndian.Uint64(out)
	var headSize = int(unsafe.Sizeof(usize))
	return out[headSize : headSize+int(usize)], nil
}

// PCKS padding
func padKey(key []byte) []byte {
	var targetKeySize = 16
	if len(key) <= 16 {
		targetKeySize = 16
	} else if len(key) <= 24 {
		targetKeySize = 24
	} else {
		targetKeySize = 32
	}
	if r := targetKeySize - len(key); r > 0 {
		nkey := make([]byte, targetKeySize)
		copy(nkey, key)
		copy(nkey[len(key):], bytes.Repeat([]byte{byte(r)}, r))
		key = nkey
	}
	return key[:targetKeySize]
}

func doHash(input string) []byte {
	hash := sha256.New()
	io.WriteString(hash, input)
	return hash.Sum(nil)
}
