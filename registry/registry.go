package registry

import "context"

// Registry 服务插件接口
type Registry interface {
	// 插件的名称
	Name() string
	// 初始化
	Init(ctx context.Context, opts ...Option) (err error)
	// 注册服务
	Register(ctx context.Context, service *Service) (err error)
	// 卸载服务
	UnRegister(ctx context.Context, service *Service) (err error)
	// 服务发现：通过服务的名字获取服务位置信息(服务节点)
	GetService(ctx context.Context, name string) (service *Service, err error)
}
