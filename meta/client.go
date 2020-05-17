package meta

import (
	"context"

	"koala/registry"

	"google.golang.org/grpc"
)

// ClientRPCMeta rpc客户端元信息, 用于参数传递
// 字段
//   ServiceName: 服务名
//   Method: 方法名
//   Nodes: 服务节点
//   CurNode: 当前选中的节点
//   Conn: 与服务的链接
type ClientRPCMeta struct {
	ServiceName string
	Method      string
	Nodes       []*registry.Node
	CurNode     *registry.Node
	HistoryNode []*registry.Node
	Conn        *grpc.ClientConn
}

// ClientOption 修改ClientRPCMeta字段
type ClientOption func(m *ClientRPCMeta)

// WithClientServiceName 修改ClientRPCMeta.ServiceName
func WithClientServiceName(name string) ClientOption {
	return func(m *ClientRPCMeta) {
		m.ServiceName = name
	}
}

// WithClientMethod 修改ClientRPCMeta.Method
func WithClientMethod(method string) ClientOption {
	return func(m *ClientRPCMeta) {
		m.Method = method
	}
}

type clientRPCMetaKey struct {
}

// InitClientRPCMeta 初始化元信息
func InitClientRPCMeta(ctx context.Context, opts ...ClientOption) context.Context {
	meta := &ClientRPCMeta{}
	for _, opt := range opts {
		opt(meta)
	}
	ctx = context.WithValue(ctx, clientRPCMetaKey{}, meta)
	return ctx
}

// GetClientRPCMeta 获取元信息
func GetClientRPCMeta(ctx context.Context) *ClientRPCMeta {
	metaV := ctx.Value(clientRPCMetaKey{})
	meta, ok := metaV.(*ClientRPCMeta)
	if !ok {
		meta = &ClientRPCMeta{}
	}
	return meta
}
