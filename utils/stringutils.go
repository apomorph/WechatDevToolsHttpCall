package utils

import (
	"math/rand"
	"strconv"
	"time"
	"unicode"
)

// RandomString 随机字符串
func RandomString(le int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < le; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

// string转int
func StrToInt(s string, def int) (i int) {
	i, err := strconv.Atoi(s)
	if err != nil {
		i = def
	}
	return
}

// string转int64
func StrToInt64(s string, def int64) (i int64) {
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		i = def
	}
	return
}

// string转float32
func StrToFloat32(s string, def float32) (f float64) {
	f, err := strconv.ParseFloat(s, 32)
	if err != nil {
		f = float64(def)
	}
	return
}

// string转float64
func StrToFloat64(s string, def float64) (f float64) {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		f = def
	}
	return
}

// 字符串首字母转大写
func Ucfirst(str string) string {
	for i, v := range str {
		return string(unicode.ToUpper(v)) + str[i+1:]
	}
	return str
}

// 字符串首字母转小写
func Lcfirst(str string) string {
	for i, v := range str {
		return string(unicode.ToLower(v)) + str[i+1:]
	}
	return str
}
