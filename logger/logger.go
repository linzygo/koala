package logger

import "context"

var logMgr *LogMgr

// Start 启动日志记录
// 参数
//   opts: Option不定参
func Start(opts ...Option) {
	for _, opt := range opts {
		opt(config)
	}

	var writers []Writer
	if config.Type&LogTypeConsole != 0 {
		writers = append(writers, NewConsoleWriter())
	}
	if config.Type&LogTypeFile != 0 {
		writers = append(writers, NewFileWriter())
	}

	logMgr = &LogMgr{
		datachan: make(chan *LogData, config.ChanSize),
		writers:  writers,
	}

	go logMgr.run()
}

// Stop 停止日志记录
func Stop() {
	logMgr.stop()
}

// Debug 调试日志
func Debug(ctx context.Context, format string, args ...interface{}) {
	logMgr.printLog(ctx, LogLevelDebug, 2, format, args...)
}

// Trace 跟踪日志
func Trace(ctx context.Context, format string, args ...interface{}) {
	logMgr.printLog(ctx, LogLevelTrace, 2, format, args...)
}

// Info 信息日志
func Info(ctx context.Context, format string, args ...interface{}) {
	logMgr.printLog(ctx, LogLevelInfo, 2, format, args...)
}

// Access 请求日志
func Access(ctx context.Context, format string, args ...interface{}) {
	logMgr.printLog(ctx, LogLevelAccess, 2, format, args...)
}

// Warn 警告日志
func Warn(ctx context.Context, format string, args ...interface{}) {
	logMgr.printLog(ctx, LogLevelWarn, 2, format, args...)
}

// Error 错误日志
func Error(ctx context.Context, format string, args ...interface{}) {
	logMgr.printLog(ctx, LogLevelError, 2, format, args...)
}
