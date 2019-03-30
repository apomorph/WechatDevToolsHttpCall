package utils

import (
	"strconv"
	"time"
)

// GetTimestamp 当前时间戳
func GetTimestamp() string {
	return strconv.FormatInt(time.Now().Unix(), 10)
}
