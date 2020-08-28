package crypto

import (
	"crypto/sha256"
	"encoding/hex"
)

func Hash(msg string) string {
	hashFunc := sha256.New()
	hashFunc.Write([]byte(msg))
	ret := hashFunc.Sum(nil)
	return hex.EncodeToString(ret)
}
