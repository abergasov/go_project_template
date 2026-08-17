package utils

import (
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"hash"
	"sync"
)

var (
	sha256Pool = sync.Pool{New: func() any {
		return sha256.New()
	}}

	sha512Pool = sync.Pool{New: func() any {
		return sha512.New()
	}}
)

func calcSHA256hash(data []byte) []byte {
	h, ok := sha256Pool.Get().(hash.Hash)
	if !ok {
		h = sha256.New()
	}
	defer sha256Pool.Put(h)
	h.Reset()

	// The hash interface never returns an error, for that reason
	// we are not handling the error below. For reference, it is
	// stated here https://golang.org/pkg/hash/#Hash

	// #nosec G104
	_, _ = h.Write(data)
	var sum [sha256.Size]byte
	h.Sum(sum[:0])
	return sum[:]
}

func calcSHA512hash(data []byte) []byte {
	h, ok := sha512Pool.Get().(hash.Hash)
	if !ok {
		h = sha512.New()
	}
	defer sha512Pool.Put(h)
	h.Reset()

	// The hash interface never returns an error, for that reason
	// we are not handling the error below. For reference, it is
	// stated here https://golang.org/pkg/hash/#Hash

	// #nosec G104
	_, _ = h.Write(data)
	var sum [sha512.Size]byte
	h.Sum(sum[:0])
	return sum[:]
}

// HashSHA256 computes SHA-256 hash returns hex-encoded string
func HashSHA256(data []byte) string {
	return hex.EncodeToString(calcSHA256hash(data))
}

// HashSHA512 computes SHA-512 hash, returns hex-encoded string
func HashSHA512(data []byte) string {
	return hex.EncodeToString(calcSHA512hash(data))
}
