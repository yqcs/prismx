package cryptoPlus

import (
	"crypto/md5"
	"fmt"
	"io"

	"github.com/yqcs/prismx/utils/logger"
)

func ToMD5(str string) string {
	m := md5.New()
	_, err := io.WriteString(m, str)
	if err != nil {
		logger.Error(err.Error())
	}
	arr := m.Sum(nil)
	return fmt.Sprintf("%x", arr)
}
