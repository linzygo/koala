package etcd

import (
	"context"
	"fmt"
	"koala/registry"
	"testing"
	"time"
)

func TestRegistry(t *testing.T) {
	ctx, cancle := context.WithTimeout(context.TODO(), time.Second*60)
	defer cancle()
	registryInst, err := registry.InitPlugin(ctx, "etcd",
		registry.WithAddrs([]string{"118.31.1.215:2379"}),
		registry.WithHeartBeat(5),
		registry.WithRegistryPath("/com.lzy.golang/koala"),
		registry.WithTimeout(time.Second),
	)
	if err != nil {
		t.Errorf("初始化注册插件失败, err=%v", err)
		return
	}
	service := &registry.Service{
		Name: "test_service",
	}

	service.Nodes = append(service.Nodes,
		&registry.Node{
			IP:   "127.0.0.1",
			Port: 8801,
		},
		&registry.Node{
			IP:   "127.0.0.2",
			Port: 8801,
		},
	)

	err = registryInst.Register(ctx, service)
	if err != nil {
		t.Errorf("注册服务失败, err=%v", err)
		return
	}
	go func() {
		time.Sleep(time.Second * 10)
		service.Nodes = append(service.Nodes, &registry.Node{
			IP:   "127.0.0.3",
			Port: 8801,
		},
			&registry.Node{
				IP:   "127.0.0.4",
				Port: 8801,
			},
		)
	}()
	fmt.Println("注册成功, 开始检查获取服务")
	for {
		service, err := registryInst.GetService(ctx, "test_service")
		if err != nil {
			t.Errorf("获取服务失败, err=%v", err)
			return
		}

		for _, node := range service.Nodes {
			fmt.Printf("service:%s, node:%#v\n", service.Name, node)
		}
		fmt.Println()
		time.Sleep(time.Second * 5)
	}
}
