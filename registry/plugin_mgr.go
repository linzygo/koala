package registry

import (
	"context"
	"fmt"
	"sync"
)

var (
	pluginMgr = &PluginMgr{
		plugins: make(map[string]Registry),
	}
)

// PluginMgr 注册插件管理器
type PluginMgr struct {
	plugins map[string]Registry
	lock    sync.Mutex
}

func (p *PluginMgr) registerPlugin(plugin Registry) (err error) {
	p.lock.Lock()
	defer p.lock.Unlock()

	if _, ok := p.plugins[plugin.Name()]; ok {
		err = fmt.Errorf("%s插件不能重复注册", plugin.Name())
		return
	}

	p.plugins[plugin.Name()] = plugin
	return
}

func (p *PluginMgr) initPlugin(ctx context.Context, name string, opts ...Option) (plugin Registry, err error) {
	p.lock.Lock()
	defer p.lock.Unlock()

	plugin, ok := p.plugins[name]

	if !ok {
		err = fmt.Errorf("%s插件不存在，不能初始化", name)
		return
	}

	plugin.Init(ctx, opts...)

	return
}

// RegisterPlugin 注册插件
func RegisterPlugin(plugin Registry) error {
	return pluginMgr.registerPlugin(plugin)
}

// InitPlugin 初始化插件
func InitPlugin(ctx context.Context, name string, opts ...Option) (Registry, error) {
	return pluginMgr.initPlugin(ctx, name, opts...)
}
