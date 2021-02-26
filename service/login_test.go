package service

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
	"xg/crypto"
)

const charset = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ!@#$%^&*()_"

func generateOneRandChar() uint8{
	rand := rand.Uint32()
	offset := int(rand) % len(charset)
	passwordPart := charset[offset]
	return passwordPart
}
func b2s(bs []uint8) string {
	ba := []byte{}
	for _, b := range bs {
		ba = append(ba, b)
	}
	return string(ba)
}
func TestGeneratePassword(t *testing.T) {
	rand.Seed(time.Now().Unix() + 50)
	passwordSize := 16
	passwordBytes := make([]byte, passwordSize)
	for i := 0; i < passwordSize; i ++ {
		passwordBytes[i] = generateOneRandChar()
	}
	passwordStr := b2s(passwordBytes)
	fmt.Println(passwordStr)

	fmt.Println(crypto.Hash(passwordStr))
}
