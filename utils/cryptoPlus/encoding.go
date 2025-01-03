package cryptoPlus

import (
	"bytes"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"html"
	"io"
	"regexp"
	"strings"
)

var (
	charsets = []string{"utf-8", "gbk", "gb2312", "latin1", "iso-8859-1"}
)

// TransCode 转码
func TransCode(body []byte, encode string) string {
	if strings.Contains(strings.ToLower(encode), "gb") {
		O := transform.NewReader(bytes.NewReader(body), simplifiedchinese.GBK.NewDecoder())
		decoder, err := io.ReadAll(O)
		if err == nil {
			body = decoder
		}
	} else if strings.Contains(strings.ToLower(encode), "iso-8859-1") || strings.Contains(strings.ToLower(encode), "latin1") {
		decoder, _, err := transform.Bytes(charmap.Windows1252.NewEncoder(), body)
		if err == nil {
			body = decoder
		}
	}
	return html.UnescapeString(string(body))
}

// GetEncoding 获取编码
func GetEncoding(contentType string, body []byte) string {
	r1, err := regexp.Compile(`(?im)charset=\s*?([\w-]+)`)
	if err != nil {
		return ""
	}
	headerCharset := r1.FindString(contentType)
	if headerCharset != "" {
		for _, v := range charsets {
			if strings.Contains(strings.ToLower(headerCharset), v) {
				return v
			}
		}
	}

	r2, err := regexp.Compile(`(?im)<meta.*?charset=['"]?([\w-]+)["']?.*?>`)
	if err != nil {
		return ""
	}
	htmlCharset := r2.FindString(string(body))
	if htmlCharset != "" {
		for _, v := range charsets {
			if strings.Contains(strings.ToLower(htmlCharset), v) {
				return v
			}
		}
	}
	return "utf-8"
}
