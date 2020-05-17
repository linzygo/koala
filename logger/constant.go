package logger

import "strings"

// LogType 日志类型
type LogType int

/*
 * LogType 值
 */
const (
	// 控制台
	LogTypeConsole = 1 << iota
	// 文件
	LogTypeFile

	// 对应的字符串描述
	LogTypeConsoleStr = "console"
	LogTypeFileStr    = "file"
)

func (lt LogType) String() string {
	switch lt {
	case LogTypeFile:
		return LogTypeFileStr
	case LogTypeConsole:
		return LogTypeConsoleStr
	default:
		return "unknow"
	}
}

func getLogType(typeStr string) LogType {
	typeStr = strings.ToLower(typeStr)
	logType := 0
	if strings.Index(typeStr, LogTypeFileStr) != -1 {
		logType |= LogTypeFile
	}
	if strings.Index(typeStr, LogTypeConsoleStr) != -1 {
		logType |= LogTypeConsole
	}
	if logType == 0 {
		logType = LogTypeConsole
	}
	return LogType(logType)
}

// LogLevel 日志级别
type LogLevel int

/*
 * LogLevel 值
 */
const (
	_ = iota
	// 调试日志, 最低级别
	LogLevelDebug
	// 跟踪日志
	LogLevelTrace
	// 信息日志
	LogLevelInfo
	// 请求访问日志
	LogLevelAccess
	// 警告日志，可能发生异常的地方
	LogLevelWarn
	// 错误日志，程序可以正常运行
	LogLevelError

	// 对应的字符串描述
	LogLevelDebugStr  = "DEBUG"
	LogLevelTraceStr  = "TRACE"
	LogLevelInfoStr   = "INFO"
	LogLevelAccessStr = "ACCESS"
	LogLevelWarnStr   = "WARN"
	LogLevelErrorStr  = "ERROR"
)

func (level LogLevel) String() string {
	switch level {
	case LogLevelDebug:
		return LogLevelDebugStr
	case LogLevelTrace:
		return LogLevelTraceStr
	case LogLevelInfo:
		return LogLevelInfoStr
	case LogLevelAccess:
		return LogLevelAccessStr
	case LogLevelWarn:
		return LogLevelWarnStr
	case LogLevelError:
		return LogLevelErrorStr
	default:
		return "UNKNOW"
	}
}

func getLogLevel(levelStr string) LogLevel {
	switch strings.ToUpper(levelStr) {
	case LogLevelDebugStr:
		return LogLevelDebug
	case LogLevelTraceStr:
		return LogLevelTrace
	case LogLevelInfoStr:
		return LogLevelInfo
	case LogLevelAccessStr:
		return LogLevelAccess
	case LogLevelWarnStr:
		return LogLevelWarn
	case LogLevelErrorStr:
		return LogLevelError
	default:
		return LogLevelInfo
	}
}

// Color 日志级别对应的颜色
func (level LogLevel) Color() Color {
	switch level {
	case LogLevelDebug:
		return White
	case LogLevelTrace:
		return Yellow
	case LogLevelInfo:
		return Green
	case LogLevelAccess:
		return Blue
	case LogLevelWarn:
		return Cyan
	case LogLevelError:
		return Red
	}
	return Magenta
}
