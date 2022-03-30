package utils

import (
	"crypto/sha1"
	"encoding/hex"
)

// Sha1 encode bytes to sha1 hex string
func Sha1(data []byte) string {
	h := sha1.New()
	h.Write(data)
	bs := h.Sum(nil)
	return hex.EncodeToString(bs)
}
