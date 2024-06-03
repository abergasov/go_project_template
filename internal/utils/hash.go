package utils

import (
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
)

func HashSHA256(data []byte) string {
	hash := sha256.New()
	hash.Write(data)
	hashBytes := hash.Sum(nil)
	return hex.EncodeToString(hashBytes)
}

func HashSHA512(data []byte) string {
	hash := sha512.New()
	hash.Write(data)
	hashBytes := hash.Sum(nil)
	return hex.EncodeToString(hashBytes)
}
