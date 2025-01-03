package arr

import (
	"bytes"
	"golang.org/x/net/html/charset"
	"io"
	"net/url"
	"strings"
)

// DeleteSliceValueToLower 删除切片里指定的值
func DeleteSliceValueToLower(list []string, value string) []string {
	for i := 0; i < len(list); i++ {
		if strings.ToLower(list[i]) == strings.ToLower(value) {
			list = append(list[:i], list[i+1:]...)
		}
	}
	return list
}

// SliceRemoveDuplicates 去除string数组重复项
func SliceRemoveDuplicates(slice []string) (result []string) {
	temp := map[string]struct{}{}
	for _, item := range slice {
		if _, ok := temp[item]; !ok {
			temp[item] = struct{}{}
			result = append(result, item)
		}
	}
	return result
}

// IntSliceRemoveDuplicates 去除int数组重复项
func IntSliceRemoveDuplicates(slice []int) (result []int) {
	temp := map[int]struct{}{}
	for _, item := range slice {
		if _, ok := temp[item]; !ok {
			temp[item] = struct{}{}
			result = append(result, item)
		}
	}
	return result
}

// DeleteSliceValue 删除切片里指定的值
func DeleteSliceValue(list []string, value string) []string {
	for i := 0; i < len(list); i++ {
		if list[i] == value {
			list = append(list[:i], list[i+1:]...)
		}
	}
	return list
}

// IntSliceContains int数组是否包含指定数据
func IntSliceContains(sl []int, v int) bool {
	for _, vv := range sl {
		if vv == v {
			return true
		}
	}
	return false
}

// IsContain 数组是否包含指定数据
func IsContain(sl []string, v string) bool {
	for _, vv := range sl {
		if vv == v {
			return true
		}
	}
	return false
}
func IsContainByUrl(sl []string, v string) bool {
	for _, vv := range sl {
		u1, err := url.Parse(vv)
		if err != nil {
			return false
		}
		u2, err := url.Parse(v)
		if err != nil {
			return false
		}
		if u1.Host == u2.Host {
			return true
		}
		if vv == v {
			return true
		}
	}
	return false
}

// IsContainToLower 不区分大小写
func IsContainToLower(sl []string, v string) bool {
	for _, vv := range sl {
		if strings.ToLower(vv) == strings.ToLower(v) {
			return true
		}
	}
	return false
}

// FuzzContainToLower 数组模糊包含
func FuzzContainToLower(sl []string, v string) bool {
	for _, vv := range sl {
		if strings.Contains(strings.ToLower(vv), strings.ToLower(v)) {
			return true
		}
	}
	return false
}

func ConvResponse(b []byte) string {
	var r1 []rune
	for _, i := range b {
		r1 = append(r1, rune(i))
	}
	return string(r1)
}

func ConvertUTF8(c io.Reader, contentType string) (bytes.Buffer, error) {
	buff := bytes.Buffer{}
	content, err := charset.NewReader(c, contentType)
	if err != nil {
		return buff, err
	}
	_, err = io.Copy(&buff, content)
	if err != nil {
		return buff, err
	}
	return buff, nil
}
