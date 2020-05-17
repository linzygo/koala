package logger

import (
	"testing"
	"time"
)

func TestFileWriter(t *testing.T) {
	writer := NewFileWriter()
	filename, funcname, line := GetCallStack(1)
	logData := &LogData{
		msg:      "测试呀测试呀",
		timeStr:  time.Now().Format("2006-01-02 15:04:05.000"),
		filename: filename,
		funcname: funcname,
		lineno:   line,
		loglevel: LogLevelAccess,
	}
	writer.Write(logData)
	logData = &LogData{
		msg:      "测试呀测试呀",
		timeStr:  time.Now().Format("2006-01-02 15:04:05.000"),
		filename: filename,
		funcname: funcname,
		lineno:   line,
		loglevel: LogLevelInfo,
	}
	writer.Write(logData)
	logData = &LogData{
		msg:      "测试呀测试呀",
		timeStr:  time.Now().Format("2006-01-02 15:04:05.000"),
		filename: filename,
		funcname: funcname,
		lineno:   line,
		loglevel: LogLevelWarn,
	}
	writer.Write(logData)
	writer.Close()
}
