package middleware

import (
	"context"
	"fmt"
	"koala/errno"
	"koala/logger"
	"koala/meta"

	"google.golang.org/grpc"
)

// RPCShortConnectMiddleware 建立rpc短链接中间件
func RPCShortConnectMiddleware(handle HandleFunc) HandleFunc {
	return func(ctx context.Context, req interface{}) (resp interface{}, err error) {
		rpcMeta := meta.GetClientRPCMeta(ctx)
		if rpcMeta.CurNode == nil {
			logger.Error(ctx, "当前节点不是有效节点")
			err = errno.InvalidNode
			return
		}
		logger.Debug(ctx, "rpc[%s.%s]当前使用的节点%#v", rpcMeta.ServiceName, rpcMeta.Method, rpcMeta.CurNode)
		addr := fmt.Sprintf("%s:%d", rpcMeta.CurNode.IP, rpcMeta.CurNode.Port)
		conn, err := grpc.Dial(addr, grpc.WithInsecure())
		if err != nil {
			logger.Error(ctx, "连接服务[%s.%s][%s]失败, err=%v", rpcMeta.ServiceName, rpcMeta.Method, addr, err)
			return nil, errno.ConnectFail
		}

		logger.Debug(ctx, "连接成功")
		rpcMeta.Conn = conn
		defer func() {
			conn.Close()
			rpcMeta.Conn = nil
		}()
		resp, err = handle(ctx, req)

		return
	}
}
