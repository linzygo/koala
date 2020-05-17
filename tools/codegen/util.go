package main

import (
	"bytes"
	"fmt"
	"os"
	"text/template"
	"unicode"
)

// ToUnderscoreString 把驼峰格式的字符串转成全小写下划线风格
func ToUnderscoreString(src string) string {
	var buf bytes.Buffer
	for index, c := range []rune(src) {
		if unicode.IsUpper(c) {
			if index != 0 {
				buf.WriteByte('_')
			}
			buf.WriteRune(unicode.ToLower(c))
		} else {
			buf.WriteRune(c)
		}
	}
	return buf.String()
}

// SaveCodeToFile 根据模板生成代码并保存到指定文件
// 参数
//   fpath: 文件路径
//   templateName: 模板名称, 用于错误发生时输出日志, 明确是生成哪个代码时发生错误
//   templateStr: 模板
//   metaData: 用于替换模板参数的数据
// 返回值
//   err: error
func SaveCodeToFile(fpath, templateName, templateStr string, metaData interface{}) (err error) {
	file, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		fmt.Printf("创建文件[%s]失败, err=%v\n", fpath, err)
		return
	}
	defer file.Close()

	t := template.New("main")
	t, err = t.Parse(templateStr)
	if err != nil {
		fmt.Printf("%s模板解析失败, err=%v\n", templateName, err)
		return
	}

	err = t.Execute(file, metaData)
	if err != nil {
		fmt.Printf("生成%s代码失败, err=%v\n", templateName, err)
	}

	return
}
