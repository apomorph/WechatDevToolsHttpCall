package utils

import (
	"fmt"
	"log"
	"runtime"
)

// Debug 调试日志
func Debug(v ...interface{}) {
	logPrint("[DEBUG] ", v...)
}

// Info 信息
func Info(v ...interface{}) {
	logPrint("[INFO] ", v...)
}

// Error 错误
func Error(v ...interface{}) {
	logPrint("[ERROR] ", v...)
}

func logPrint(prefix string, v ...interface{}) {
	pc, _, line, _ := runtime.Caller(2)
	f := runtime.FuncForPC(pc)
	log.SetPrefix(prefix)
	log.Printf("%v:%v %v\n", f.Name(), line, fmt.Sprint(v...))
}
