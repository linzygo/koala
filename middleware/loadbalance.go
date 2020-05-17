package middleware

import (
	"context"
	"koala/errno"
	"koala/loadbalance"
	"koala/logger"
	"koala/meta"
)

// NewLoadBalanceMiddleware 创建负载均衡中间件
func NewLoadBalanceMiddleware(balancer loadbalance.LoadBalance) Middleware {
	return func(handle HandleFunc) HandleFunc {
		return func(ctx context.Context, req interface{}) (resp interface{}, err error) {
			rpcMeta := meta.GetClientRPCMeta(ctx)
			// 构造用于传递已选过的节点的context
			ctx = loadbalance.WithSelectedNodes(ctx)
			for {
				rpcMeta.CurNode, err = balancer.Select(ctx, rpcMeta.Nodes)
				if err != nil {
					logger.Error(ctx, "rpc[%s.%s]选择节点失败:%#v", rpcMeta.ServiceName, rpcMeta.Method, rpcMeta.CurNode)
					return
				}
				rpcMeta.HistoryNode = append(rpcMeta.HistoryNode, rpcMeta.CurNode)
				logger.Debug(ctx, "rpc[%s.%s]选择的节点:%#v", rpcMeta.ServiceName, rpcMeta.Method, rpcMeta.CurNode)
				resp, err = handle(ctx, req)
				logger.Debug(ctx, "handle rpc[%s.%s] resp=%#v, err=%v", rpcMeta.ServiceName, rpcMeta.Method, resp, err)
				if err != nil {
					if errno.IsConnectFail(err) {
						continue
					}
					return
				}
				// 成功
				break
			}
			return
		}
	}
}
