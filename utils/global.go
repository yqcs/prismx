package utils

import (
	"prismx_cli/utils/arr"
)

const (
	DateFormat string = "2006-01-02 15:04:05"
)

// GlobalError 通用检查错误信息
func GlobalError(err error) bool {
	if err == nil {
		return false
	}
	errs := []string{
		"closed by the remote host", "too many connections",
		"i/o timeout", "A connection attempt failed",
		"established connection failed", "connection attempt failed",
		"Unable to read", "is not allowed to connect to this",
		"no pg_hba.conf entry",
		"An existing connection was forcibly closed by the remote host",
		"No connection could be made",
		"local file '/etc/group' is not registered",
		"unexpected EOF",
		"invalid packet size",
		"bad connection",
	}
	return arr.IsContain(errs, err.Error())
}
