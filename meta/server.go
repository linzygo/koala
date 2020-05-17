package meta

import "context"

// ServerMeta rpc服务端元信息, 用于打点采样
// 字段:
//   ServiceName 服务名;
//   Method 方法名;
//   Cluster 集群;
//   TraceID 请求的trace id;
//   ServerIP 当前服务器IP;
//   ClientIP 请求的客户端IP;
//   IDC 机房;
type ServerMeta struct {
	ServiceName string
	Method      string
	Cluster     string
	TraceID     string
	ServerIP    string
	ClientIP    string
	IDC         string
}

// ServerOption 修改ServerMeta字段
type ServerOption func(sm *ServerMeta)

// WithServerServiceName 修改ServerMeta.ServiceName
func WithServerServiceName(name string) ServerOption {
	return func(sm *ServerMeta) {
		sm.ServiceName = name
	}
}

// WithServerMethod 修改ServerMeta.Method
func WithServerMethod(method string) ServerOption {
	return func(sm *ServerMeta) {
		sm.Method = method
	}
}

// WithServerCluster 修改ServerMeta.Cluster
func WithServerCluster(cluster string) ServerOption {
	return func(sm *ServerMeta) {
		sm.Cluster = cluster
	}
}

// WithServerTraceID 修改ServerMeta.TraceID
func WithServerTraceID(traceID string) ServerOption {
	return func(sm *ServerMeta) {
		sm.TraceID = traceID
	}
}

// WithServerServerIP 修改ServerMeta.ServerIP
func WithServerServerIP(ip string) ServerOption {
	return func(sm *ServerMeta) {
		sm.ServerIP = ip
	}
}

// WithServerClientIP 修改ServerMeta.ClientIP
func WithServerClientIP(ip string) ServerOption {
	return func(sm *ServerMeta) {
		sm.ClientIP = ip
	}
}

// WithServerIDC 修改ServerMeta.IDC
func WithServerIDC(idc string) ServerOption {
	return func(sm *ServerMeta) {
		sm.IDC = idc
	}
}

type serverMetaContextKey struct{}

// GetServerMeta 获取服务器元信息
func GetServerMeta(ctx context.Context) *ServerMeta {
	meta, ok := ctx.Value(serverMetaContextKey{}).(*ServerMeta)
	if !ok {
		meta = &ServerMeta{}
	}
	return meta
}

// InitServerMeta 初始化服务器元信息
func InitServerMeta(ctx context.Context, opts ...ServerOption) context.Context {
	meta := &ServerMeta{}
	for _, opt := range opts {
		opt(meta)
	}
	return context.WithValue(ctx, serverMetaContextKey{}, meta)
}
