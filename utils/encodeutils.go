package utils

import (
	"encoding/base64"
	"net/url"
)

// url encode
func UrlEncode(param string) (result string) {
	return url.QueryEscape(param)
}

// Base46Encode base64加密
func Base46Encode(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}

// Base46Decode base64解密
func Base46Decode(str string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(str)
}
