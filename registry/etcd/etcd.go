package etcd

import (
	"context"
	"encoding/json"
	"fmt"
	"koala/registry"
	"path"
	"sync"
	"sync/atomic"
	"time"

	"go.etcd.io/etcd/clientv3"
)

const (
	// InitServiceNum 服务数量初始化值
	InitServiceNum = 8
	// SyncServiceInterval 服务发现获取服务信息的时间间隔
	SyncServiceInterval = time.Second * 10
)

// Registry etcd注册插件
type Registry struct {
	cfg                *registry.Config
	client             *clientv3.Client
	serviceCh          chan *registry.Service // 用于把服务注册转到goroutine完成
	registryServiceMap map[string]*RegistryService
	serviceCache       atomic.Value // 用于保存服务缓存, 原子操作从而避免加锁，取出来后只能遍历不能修改，避免并发冲突
	lock               sync.Mutex   // 用于加锁阻塞，避免多个请求同时并发，把etcd压垮
}

// RegistryService 服务与etcd的关联
type RegistryService struct {
	id          clientv3.LeaseID
	service     *registry.Service
	isRegister  bool // 用于标志该服务是否已经注册到etcd
	keepAliveCh <-chan *clientv3.LeaseKeepAliveResponse
}

// AllServiceInfo 作为缓存，用于保存etcd上的服务
type AllServiceInfo struct {
	serviceMap map[string]*registry.Service
}

var (
	etcdRegistry = &Registry{
		serviceCh:          make(chan *registry.Service, InitServiceNum),
		registryServiceMap: make(map[string]*RegistryService, InitServiceNum),
	}
)

func init() {
	allServiceInfo := &AllServiceInfo{
		serviceMap: make(map[string]*registry.Service, InitServiceNum),
	}
	etcdRegistry.serviceCache.Store(allServiceInfo)

	registry.RegisterPlugin(etcdRegistry) // 注册插件

	go etcdRegistry.run()
}

// Name 实现registry.Registry接口
func (r *Registry) Name() string {
	return "etcd"
}

// Init 实现registry.Registry接口
func (r *Registry) Init(ctx context.Context, opts ...registry.Option) (err error) {
	r.cfg = &registry.Config{}
	for _, op := range opts {
		op(r.cfg)
	}

	r.client, err = clientv3.New(clientv3.Config{
		Endpoints:   r.cfg.Addrs,
		DialTimeout: r.cfg.Timeout,
	})

	if err != nil {
		err = fmt.Errorf("初始化etcd注册插件失败, err=%v", err)
	}

	return
}

// Register 实现registry.Register接口
func (r *Registry) Register(ctx context.Context, service *registry.Service) (err error) {
	// 注册信息需要写到etcd，不能直接在这里操作，避免阻塞，需要扔到goroutine完成，同时防止chan缓冲区满了
	select {
	case r.serviceCh <- service:
	default:
		err = fmt.Errorf("etcd注册组件的service chan已经满了")
	}
	return
}

// UnRegister 实现registry.Register接口
func (r *Registry) UnRegister(ctx context.Context, service *registry.Service) (err error) {
	return
}

// GetService 实现registry.Register接口
func (r *Registry) GetService(ctx context.Context, name string) (service *registry.Service, err error) {
	// 先从缓冲区取
	service, ok := r.getServiceFromCache(name)
	if ok {
		return
	}

	// 缓冲区没有, 从etcd中拉取, 先加锁避免高并发压垮etcd，得到锁后先检查是否已经拉取下来了
	r.lock.Lock()
	defer r.lock.Unlock()
	service, ok = r.getServiceFromCache(name)
	if ok {
		return
	}

	// 从etcd中获取指定的服务信息
	key := r.servicePath(name)
	// 获取前缀包含key的信息
	resp, err := r.client.Get(ctx, key, clientv3.WithPrefix())
	if err != nil {
		err = fmt.Errorf("没有找到服务%s", name)
		return
	}

	service = &registry.Service{
		Name: name,
	}
	for _, kvs := range resp.Kvs {
		node := &registry.Node{}
		err = json.Unmarshal(kvs.Value, node)
		if err != nil {
			return
		}
		service.Nodes = append(service.Nodes, node)
	}
	// 保存到缓冲区
	allServiceInfo := r.serviceCache.Load().(*AllServiceInfo)
	allServiceInfoNew := &AllServiceInfo{
		serviceMap: make(map[string]*registry.Service, len(allServiceInfo.serviceMap)+1),
	}
	for k, v := range allServiceInfo.serviceMap {
		allServiceInfoNew.serviceMap[k] = v
	}

	allServiceInfoNew.serviceMap[name] = service
	r.serviceCache.Store(allServiceInfoNew)
	return
}

func (r *Registry) run() {
	ticker := time.NewTicker(time.Second * 10)
	for {
		select {
		case service := <-r.serviceCh:
			// 连接etcd，注册服务
			registryService, ok := r.registryServiceMap[service.Name]
			if ok {
				// 这里可以优化，同样的服务，同样的节点，就没必要增加进来了，时间关系暂不考虑了
				registryService.service.Nodes = append(registryService.service.Nodes, service.Nodes...)
				registryService.isRegister = false // 标志为未注册，使服务信息添加到etcd
				break
			}
			registryService = &RegistryService{
				service: service,
			}
			r.registryServiceMap[service.Name] = registryService
		case <-ticker.C:
			// 更新服务缓冲
			r.syncServiceFromEtcd()
		default:
			time.Sleep(time.Millisecond * 500)
			r.registerOrKeepAlive()
		}
	}
}

func (r *Registry) registerOrKeepAlive() {
	for _, registryService := range r.registryServiceMap {
		if registryService.isRegister {
			// 检查租约是否正常
			r.keepAlive(registryService)
			continue
		}
		r.registerService(registryService)
	}
}

func (r *Registry) keepAlive(registryService *RegistryService) (err error) {
	// 检查租约是否正常
	select {
	case resp := <-registryService.keepAliveCh:
		if resp == nil { // 租约不存在了
			registryService.isRegister = false
		}
	}
	return
}

func (r *Registry) registerService(registryService *RegistryService) (err error) {
	ctx, cancle := context.WithTimeout(context.TODO(), time.Second*30)
	defer cancle()
	resp, err := r.client.Grant(ctx, r.cfg.HeartBeat)
	if err != nil {
		err = fmt.Errorf("注册服务失败，err=%v", err)
		return
	}
	registryService.id = resp.ID
	for _, node := range registryService.service.Nodes {
		nodeInfo := registry.NodeInfo{
			ServiceName: registryService.service.Name,
			Node:        *node,
		}

		data, err := json.Marshal(node)
		if err != nil {
			continue
		}

		key := r.serviceNodePath(nodeInfo)
		// 注册节点, 获得租约
		_, err = r.client.Put(ctx, key, string(data), clientv3.WithLease(registryService.id))
		if err != nil {
			continue
		}

		// 设置租约永远有效
		ch, err := r.client.KeepAlive(ctx, registryService.id)
		if err != nil {
			continue
		}
		registryService.keepAliveCh = ch
		registryService.isRegister = true
	}
	return
}

func (r *Registry) serviceNodePath(nodeInfo registry.NodeInfo) string {
	addrs := fmt.Sprintf("%s:%d", nodeInfo.IP, nodeInfo.Port)
	// 使用path自动增加/
	return path.Join(r.cfg.RegistryPath, nodeInfo.ServiceName, addrs)
}

func (r *Registry) servicePath(name string) string {
	return path.Join(r.cfg.RegistryPath, name)
}

func (r *Registry) getServiceFromCache(name string) (service *registry.Service, ok bool) {
	allServiceInfo := r.serviceCache.Load().(*AllServiceInfo)
	service, ok = allServiceInfo.serviceMap[name]
	return
}

func (r *Registry) syncServiceFromEtcd() {
	allServiceInfoNew := &AllServiceInfo{
		serviceMap: make(map[string]*registry.Service, InitServiceNum),
	}

	ctx, cancle := context.WithTimeout(context.TODO(), time.Second*30)
	defer cancle()
	allServiceInfo := r.serviceCache.Load().(*AllServiceInfo)
	for _, service := range allServiceInfo.serviceMap {
		key := r.servicePath(service.Name)
		resp, err := r.client.Get(ctx, key, clientv3.WithPrefix())
		if err != nil {
			allServiceInfoNew.serviceMap[service.Name] = service
			continue
		}

		serviceNew := &registry.Service{
			Name: service.Name,
		}
		for _, kv := range resp.Kvs {
			node := &registry.Node{}
			err = json.Unmarshal(kv.Value, node)
			if err != nil {
				err = fmt.Errorf("更新服务缓存失败, err=%v", err)
				return
			}
			serviceNew.Nodes = append(serviceNew.Nodes, node)
		}
		allServiceInfoNew.serviceMap[serviceNew.Name] = serviceNew
	}

	r.serviceCache.Store(allServiceInfoNew)
}
