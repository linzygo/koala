package logger

// Writer 日志输出器接口
type Writer interface {
	Write(data *LogData)
	Close()
}
