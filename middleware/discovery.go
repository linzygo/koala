package middleware

import (
	"context"
	"koala/logger"
	"koala/meta"
	"koala/registry"
)

// NewDiscoveryMiddleware 创建服务发现中间件
func NewDiscoveryMiddleware(discovery registry.Registry) Middleware {
	return func(handle HandleFunc) HandleFunc {
		return func(ctx context.Context, req interface{}) (resp interface{}, err error) {
			rpcMeta := meta.GetClientRPCMeta(ctx)

			if len(rpcMeta.Nodes) != 0 {
				resp, err = handle(ctx, req)
				return
			}

			service, err := discovery.GetService(ctx, rpcMeta.ServiceName)
			if err != nil {
				logger.Error(ctx, "服务发现获取服务地址失败, err=%v", err)
				return
			}
			rpcMeta.Nodes = service.Nodes

			resp, err = handle(ctx, req)
			return
		}
	}
}
