package logger

import (
	"bytes"
	"fmt"
)

// LogData 日志数据
// 字段
//   msg: 用户打印的消息
//   timeStr: 日志时间
//   filename: 要打印日志的代码文件
//   funcname: 要打印日志的函数
//   lineno: 要打印日志的地方在文件第几行
//   loglevle: 要打印的日志级别
//   traceID: 分布式追踪id, 没有时则本地生成一个
//   field: 请求参数, 只用于access log
type LogData struct {
	msg      string
	timeStr  string
	filename string
	funcname string
	lineno   int
	loglevel LogLevel
	traceID  string
	field    *AccessField
}

// String LogData转成固定格式的字符串
// 返回值
//   string: 格式化后的字符串
func (data *LogData) String() string {
	var buffer bytes.Buffer
	buffer.WriteString(data.timeStr)
	buffer.WriteByte(' ')
	writeData(&buffer, data.loglevel.String())
	writeData(&buffer, data.filename, data.funcname, fmt.Sprintf("%d", data.lineno))
	writeData(&buffer, "traceid", data.traceID)

	if data.loglevel == LogLevelAccess {
		buffer.WriteByte(' ')
		writeFiled(&buffer, data.field)
	}
	buffer.WriteByte(' ')
	buffer.WriteString(data.msg)
	buffer.WriteByte('\n')

	return buffer.String()
}

func writeData(buffer *bytes.Buffer, datas ...string) {
	buffer.WriteByte('[')
	for index, data := range datas {
		if index != 0 {
			buffer.WriteByte(':')
		}
		buffer.WriteString(data)
	}
	buffer.WriteByte(']')
}

func writeFiled(buffer *bytes.Buffer, field *AccessField) {
	buffer.WriteByte('[')
	if field != nil {
		for index, kv := range field.kvs {
			if index != 0 {
				buffer.WriteByte(',')
			}
			buffer.WriteString(fmt.Sprintf("%v=%v", kv.key, kv.val))
		}
	}
	buffer.WriteByte(']')
}
