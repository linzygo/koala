package logger

import (
	"fmt"
)

// Color 控制台文字颜色
type Color uint8

/*
 * Color值
 */
const (
	Black Color = iota + 30
	Red
	Green
	Yellow
	Blue
	Magenta
	Cyan
	White
)

// Format 给文本增加颜色
// 参数
//   str: 要增加颜色的文本
// 返回值
//   string: 增加颜色后的文本
func (c Color) Format(str string) string {
	// return fmt.Sprintf("\x1b[%dm%s\x1b[0m", uint8(c), str)
	return fmt.Sprintf("\x1b[%dm%s\x1b[0m", uint8(c), str)
}
