package main

var routerTemplate = `
package router

import (
	"context"
	"koala/server"
	"koala/meta"
	"{{.PackagePrefix}}/controller"
	"{{.PackagePrefix}}/{{.PackageName}}"

	"google.golang.org/grpc/peer"
)

// ServerRouter 用于对接rpc接口
type ServerRouter struct {
}

{{range .RPCs}}
// {{.Name}} 实现接口protorpc接口
func (router *ServerRouter) {{.Name}}(ctx context.Context, req *{{$.PackageName}}.{{.RequestType}}) (
	resp *{{$.PackageName}}.{{.ReturnsType}}, err error) {
	pr, _ := peer.FromContext(ctx)
	ctx = meta.InitServerMeta(ctx,
		meta.WithServerServiceName("{{$.PackageName}}"),
		meta.WithServerMethod("{{.Name}}"),
		meta.WithServerClientIP(pr.Addr.String()))
	mw := server.BuildServerMiddleware(wm{{.Name}})
	mwResp, err := mw(ctx, req)
	if mwResp != nil {
		resp = mwResp.(*{{$.PackageName}}.{{.ReturnsType}})
	}
	return
}

func wm{{.Name}}(ctx context.Context, req interface{}) (resp interface{}, err error) {
	inst := &controller.{{.Name}}Controller{}
	err = inst.CheckParams(ctx, req.(*{{$.PackageName}}.{{.RequestType}}))
	if err != nil {
		return
	}
	resp, err = inst.Process(ctx, req.(*{{$.PackageName}}.{{.RequestType}}))
	return
}

{{end}}
`
