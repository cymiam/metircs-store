package hmac

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

func CalculateSha256Sum(data []byte, key string) string {
	h := hmac.New(sha256.New, []byte(key))
	h.Write(data)

	sum := h.Sum(nil)
	return hex.EncodeToString(sum)
}
