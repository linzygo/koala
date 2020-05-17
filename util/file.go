package util

import "os"

// PathIsExist 检查文件或文件夹是否存在
// 参数
//   path: 文件或文件夹路径
// 返回值
//   bool: true, 存在; false, 不存在
func PathIsExist(path string) bool {
	_, err := os.Stat(path)
	if err == nil || os.IsExist(err) {
		return true
	}
	return false
}
