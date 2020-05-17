package middleware

import (
	"context"
	"koala/meta"

	"github.com/afex/hystrix-go/hystrix"
)

// HystrixMiddleware 熔断中间件
func HystrixMiddleware(handle HandleFunc) HandleFunc {
	return func(ctx context.Context, req interface{}) (resp interface{}, err error) {
		rpcMeta := meta.GetClientRPCMeta(ctx)
		hystrixErr := hystrix.Do(rpcMeta.ServiceName, func() error {
			resp, err = handle(ctx, req)
			return err
		}, nil)
		if hystrixErr != nil {
			return nil, hystrixErr
		}
		return
	}
}
