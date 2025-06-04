package utils

import (
	"crypto/sha256"
	"encoding/hex"
)

func HashString(s string) string {
	hash := sha256.New()
	hash.Write([]byte(s))
	return hex.EncodeToString(hash.Sum(nil))
}
