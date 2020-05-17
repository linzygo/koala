package middleware

import (
	"bytes"
	"context"
	"fmt"
	"koala/logger"
	"koala/meta"
	"strings"
	"time"

	"google.golang.org/grpc/status"
)

// AccessLogMiddleware 请求日志中间件
func AccessLogMiddleware(handle HandleFunc) HandleFunc {
	return func(ctx context.Context, req interface{}) (resp interface{}, err error) {
		now := time.Now()
		resp, err = handle(ctx, req)

		serverMeta := meta.GetServerMeta(ctx)
		errStatus, _ := status.FromError(err)

		ctx = logger.WithFieldContext(ctx)
		logger.AddField(ctx, "cost_time", fmt.Sprintf("%dus", time.Since(now).Microseconds()))
		logger.AddField(ctx, "service", serverMeta.ServiceName)
		logger.AddField(ctx, "method", serverMeta.Method)
		logger.AddField(ctx, "server_ip", serverMeta.ServerIP)
		logger.AddField(ctx, "client_ip", serverMeta.ClientIP)
		logger.AddField(ctx, "cluster", serverMeta.Cluster)
		logger.AddField(ctx, "idc", serverMeta.IDC)
		logger.Access(ctx, "result=%v", errStatus.Code())

		return
	}
}

// RPCAccessLogMiddleware rpc请求日志中间件
func RPCAccessLogMiddleware(handle HandleFunc) HandleFunc {
	return func(ctx context.Context, req interface{}) (resp interface{}, err error) {
		now := time.Now()
		resp, err = handle(ctx, req)

		rpcMeta := meta.GetClientRPCMeta(ctx)
		errStatus, _ := status.FromError(err)

		var curNode string
		if rpcMeta.CurNode != nil {
			curNode = fmt.Sprintf("%s:%d", rpcMeta.CurNode.IP, rpcMeta.CurNode.Port)
		}
		var hisNodes bytes.Buffer
		for _, hisNode := range rpcMeta.HistoryNode {
			if hisNode != nil && hisNode != rpcMeta.CurNode {
				hisNodes.WriteString(fmt.Sprintf("%s:%d,", hisNode.IP, hisNode.Port))
			}
		}

		ctx = logger.WithFieldContext(ctx)
		logger.AddField(ctx, "cost_time", fmt.Sprintf("%dus", time.Since(now).Microseconds()))
		logger.AddField(ctx, "service", rpcMeta.ServiceName)
		logger.AddField(ctx, "method", rpcMeta.Method)
		logger.AddField(ctx, "current_node", curNode)
		logger.AddField(ctx, "history_nodes", strings.TrimSuffix(hisNodes.String(), ","))
		logger.Access(ctx, "result=%v", errStatus.Code())
		return
	}
}
