package middleware

import (
	"context"
	"koala/logger"
	"koala/util"

	"google.golang.org/grpc/metadata"
)

// PrepareMiddleware 预备中间件, 做一些全局的初始化
func PrepareMiddleware(handle HandleFunc) HandleFunc {
	return func(ctx context.Context, req interface{}) (resp interface{}, err error) {
		var traceID string
		// 从Context获取grpc的metadata
		md, ok := metadata.FromIncomingContext(ctx)
		if ok {
			vals, ok := md[util.TraceID]
			if ok && len(vals) > 0 {
				traceID = vals[0]
			}
		}

		if len(traceID) == 0 {
			traceID = logger.GenTraceID()
		}

		ctx = logger.WithTraceID(ctx, traceID)
		resp, err = handle(ctx, req)
		return
	}
}
