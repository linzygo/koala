package errno

import "fmt"

// KoalaError 自定义错误
type KoalaError struct {
	Code int
	Msg  string
}

// Error 实现error接口
func (err *KoalaError) Error() string {
	return fmt.Sprintf("{code:%d, message:\"%s\"}", err.Code, err.Msg)
}

/*
 * 自定义错误
 */
var (
	InvalidNode = &KoalaError{
		Code: 1003,
		Msg:  "无效的节点",
	}
	AllNodeFailed = &KoalaError{
		Code: 1004,
		Msg:  "所有节点失败",
	}
	EmptyNode = &KoalaError{
		Code: 1005,
		Msg:  "无节点",
	}
	ConnectFail = &KoalaError{
		Code: 1006,
		Msg:  "连接失败",
	}
)

// IsConnectFail 检查错误是否为连接错误
// 参数
//   err: 错误
// 返回值
//   bool: 如果是连接错误则返回true，否则返回false
func IsConnectFail(err error) bool {
	koalaError, ok := err.(*KoalaError)
	if ok && koalaError == ConnectFail {
		return true
	}
	return false
}
