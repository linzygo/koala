package middleware

import (
	"context"
	"koala/meta"
	"koala/middleware/prometheus"
	"time"
)

/*
 * defaultServerMetricss 服务器采样打点
 * defaultRPCMetricscs rpc客户端调用采样打点
 */
var (
	defaultServerMetrics = prometheus.NewServerMetrics()
	defaultRPCMetrics    = prometheus.NewRPCMetrics()
)

// PrometheusServerMiddleware prometheus采样打点中间件
func PrometheusServerMiddleware(next HandleFunc) HandleFunc {
	return func(ctx context.Context, req interface{}) (resp interface{}, err error) {
		serverMeta := meta.GetServerMeta(ctx)
		defaultServerMetrics.IncRequest(ctx, serverMeta.ServiceName, serverMeta.Method)
		startTime := time.Now()
		resp, err = next(ctx, req)
		us := time.Since(startTime).Microseconds()
		if err != nil {
			defaultServerMetrics.IncRequestErr(ctx, serverMeta.ServiceName, serverMeta.Method, err)
		}
		defaultServerMetrics.Cost(ctx, serverMeta.ServiceName, serverMeta.Method, us)
		return
	}
}

// PrometheusRPCMiddleware rpc客户端调用prometheus采样打点中间件
func PrometheusRPCMiddleware(next HandleFunc) HandleFunc {
	return func(ctx context.Context, req interface{}) (resp interface{}, err error) {
		rpcMeta := meta.GetClientRPCMeta(ctx)
		defaultRPCMetrics.IncRequest(ctx, rpcMeta.ServiceName, rpcMeta.Method)
		startTime := time.Now()
		resp, err = next(ctx, req)
		us := time.Since(startTime).Microseconds()
		if err != nil {
			defaultRPCMetrics.IncRequestErr(ctx, rpcMeta.ServiceName, rpcMeta.Method, err)
		}
		defaultRPCMetrics.Cost(ctx, rpcMeta.ServiceName, rpcMeta.Method, us)
		return
	}
}
