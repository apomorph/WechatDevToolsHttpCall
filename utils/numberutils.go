package utils

import (
	"math/rand"
	"strconv"
	"time"
)

// RandomNumber 随机数字
func RandomNumber(le int) string {
	str := "0123456789"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < le; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

// int转string
func IntToStr(i int) string {
	return strconv.Itoa(i) // strconv.FormatInt(int64(i), 10)
}

// int64转string
func Int64ToStr(i int64) string {
	return strconv.FormatInt(i, 10)
}
