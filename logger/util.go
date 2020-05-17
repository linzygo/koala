package logger

import (
	"path/filepath"
	"runtime"
)

// GetCallStack 获取调用函数的文件名、函数名和行号
func GetCallStack(skip int) (filename, funcname string, line int) {
	pc, file, l, ok := runtime.Caller(skip)
	if ok {
		filename = filepath.Base(file)
		line = l
		funcname = filepath.Base(runtime.FuncForPC(pc).Name())
	}
	return
}
