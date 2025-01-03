package cryptoPlus

import (
	"encoding/base64"
	"prismx_cli/utils/logger"
)

// Base64Decode 解密
func Base64Decode(str string) string {
	decoded, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		logger.Error(err.Error())
		return ""
	}
	return string(decoded)
}

// Base64Encode 加密
func Base64Encode(str string) string {
	encoded := base64.StdEncoding.EncodeToString([]byte(str))
	return encoded
}
