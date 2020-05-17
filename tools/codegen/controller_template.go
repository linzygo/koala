package main

var controllerTemplate = `
package controller

import (
	"context"
	"{{.PackagePrefix}}/{{.PackageName}}"
)

// {{.Name}}Controller 用于实现rpc {{.Name}}接口
type {{.Name}}Controller struct {
}

// CheckParams 检查参数
func (ctrl *{{.Name}}Controller) CheckParams(ctx context.Context, req *{{.PackageName}}.{{.RequestType}}) (err error) {
	return
}

// Process 处理业务
func (ctrl *{{.Name}}Controller) Process(ctx context.Context, req *{{.PackageName}}.{{.RequestType}}) (
	resp *{{.PackageName}}.{{.ReturnsType}}, err error) {
	return
}
`
