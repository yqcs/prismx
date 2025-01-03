package cryptoPlus

import (
	"golang.org/x/crypto/bcrypt"
	"prismx_cli/utils/logger"
)

// ValidateBcryptPassWd 验证密码
// 第一个参数是明文  第二个参数是密文
func ValidateBcryptPassWd(src string, passWd string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(passWd), []byte(src)); err != nil {
		logger.Error(err.Error())
		return false
	}
	return true
}

// GeneratePassWd 生成密码
func GeneratePassWd(src string) []byte {
	res, err := bcrypt.GenerateFromPassword([]byte(src), bcrypt.DefaultCost)
	if err != nil {
		logger.Error(err.Error())
	}
	return res
}
