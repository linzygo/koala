package main

var rpcClientTemplate = `
package client

import (
	"context"
	"koala/client"
	"koala/errno"
	"koala/meta"
	"{{.PackagePrefix}}/{{.PackageName}}"
)

// {{.ClientName}}Client rpc客户端
type {{.ClientName}}Client struct {
	serviceName   string
	serviceClient *client.KoalaServiceClient
}

// New{{.ClientName}}Client 新建一个rpc客户端
func New{{.ClientName}}Client() *{{.ClientName}}Client {
	c := &{{.ClientName}}Client{
		serviceName:   "{{.PackageName}}",
		serviceClient: client.NewServiceClient("{{.PackageName}}"),
	}
	return c
}
{{range .RPCs}}
// {{.Name}} 实现接口protorpc客户端接口
func (c *{{$.ClientName}}Client) {{.Name}}(ctx context.Context, req *{{$.PackageName}}.{{.RequestType}}) (
	resp *{{$.PackageName}}.{{.ReturnsType}}, err error) {
	callResp, err := c.serviceClient.Call(ctx, req, "{{.Name}}", wm{{.Name}})
	if callResp != nil {
		resp = callResp.(*{{$.PackageName}}.{{.ReturnsType}})
	}
	return
}

func wm{{.Name}}(ctx context.Context, req interface{}) (resp interface{}, err error) {
	rcpMeta := meta.GetClientRPCMeta(ctx)
	if rcpMeta.Conn == nil {
		return nil, errno.ConnectFail
	}
	client := {{$.PackageName}}.New{{$.ServiceName}}Client(rcpMeta.Conn)
	resp, err = client.{{.Name}}(ctx, req.(*{{$.PackageName}}.{{.RequestType}}))

	return
}
{{end}}
`
