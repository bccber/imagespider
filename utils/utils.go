package utils

import (
	"crypto/md5"
	"fmt"
	"os"
	"strings"
)

// 判断文件是否存在
func PathIsExist(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}

	if os.IsNotExist(err) {
		return false
	}

	return false
}

// 字符串MD5
func MD5(str string) string {
	data := []byte(strings.ToLower(str))
	hs := md5.Sum(data)
	md5Str := fmt.Sprintf("%x", hs)

	return strings.ToLower(md5Str)
}
