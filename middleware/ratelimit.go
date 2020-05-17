package middleware

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Limiter 限流器接口
type Limiter interface {
	Allow() bool
}

// NewLimiterMiddleware 创建限流器中间件
func NewLimiterMiddleware(l Limiter) Middleware {
	return func(next HandleFunc) HandleFunc {
		return func(ctx context.Context, req interface{}) (resp interface{}, err error) {
			if !l.Allow() {
				err = status.Error(codes.ResourceExhausted, "rate limited")
				return
			}
			resp, err = next(ctx, req)
			return
		}
	}
}
