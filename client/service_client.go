package client

import (
	"context"
	"koala/meta"
	"koala/middleware"
	_ "koala/registry/etcd" // 为了注册etcd插件
)

// KoalaServiceClient 封装每个服务的rpc client共用代码
// 字段
//   serviceName: 服务名称
type KoalaServiceClient struct {
	serviceName string
	mw          middleware.Middleware
}

// NewServiceClient 新建一个rpc服务客户端
func NewServiceClient(service string) *KoalaServiceClient {
	client := &KoalaServiceClient{
		serviceName: service,
	}
	client.initMiddleware()
	return client
}

// Call 封装了中间件等调用，传入要执行的函数
func (client *KoalaServiceClient) Call(ctx context.Context, req interface{}, method string, handle middleware.HandleFunc) (
	resp interface{},
	err error,
) {
	ctx = meta.InitClientRPCMeta(ctx, meta.WithClientServiceName(client.serviceName), meta.WithClientMethod(method))
	mwFunc := client.mw(handle)
	resp, err = mwFunc(ctx, req)
	return
}

func (client *KoalaServiceClient) initMiddleware() {
	mids := []middleware.Middleware{}
	mids = append(mids, middleware.RPCAccessLogMiddleware)
	if koalaConf.Trace.SwitchOn {
		mids = append(mids, middleware.TraceRPCMiddleware)
	}
	if koalaConf.Prometheus.SwitchOn {
		mids = append(mids, middleware.PrometheusRPCMiddleware)
	}
	if koalaConf.Limit.SwitchOn {
		mids = append(mids, middleware.NewLimiterMiddleware(koalaClient.limiter))
	}

	mids = append(mids, middleware.HystrixMiddleware)
	mids = append(mids, middleware.NewDiscoveryMiddleware(koalaClient.discoveryInst))
	mids = append(mids, middleware.NewLoadBalanceMiddleware(koalaClient.balancer))
	mids = append(mids, middleware.RPCShortConnectMiddleware)
	client.mw = middleware.Chain(middleware.PrepareMiddleware, mids...)
}
