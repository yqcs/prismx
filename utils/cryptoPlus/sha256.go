package cryptoPlus

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

func ToSha256(str, keyStr string) string {
	s := []byte(str)
	key := []byte(keyStr)
	m := hmac.New(sha256.New, key)
	m.Write(s)
	signature := hex.EncodeToString(m.Sum(nil))
	return signature
}

func SHA256Sum(data any) string {
	hash := sha256.New()
	if v, ok := data.([]byte); ok {
		hash.Write(v)
	} else if v, ok := data.(string); ok {
		hash.Write([]byte(v))
	} else {
		return ""
	}

	return hex.EncodeToString(hash.Sum(nil))
}
