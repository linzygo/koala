package logger

import "fmt"

// ConsoleWriter 控制台输出器
type ConsoleWriter struct {
}

// NewConsoleWriter 创建控制台输出器
func NewConsoleWriter() *ConsoleWriter {
	return &ConsoleWriter{}
}

// Write 实现接口Writer
func (writer *ConsoleWriter) Write(data *LogData) {
	color := data.loglevel.Color()
	msg := color.Format(data.String())
	fmt.Println(msg)
}

// Close 实现接口Writer
func (writer *ConsoleWriter) Close() {

}
