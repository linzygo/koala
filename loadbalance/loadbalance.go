package loadbalance

import (
	"context"
	"koala/registry"
)

/*
 * DefaultNodeWeight 节点默认权重
 * LoadBalanceTypeRandom 负载均衡类型-随机算法
 * LoadBalanceTypeRoundRobin 负载均衡类型-轮询
 */
const (
	DefaultNodeWeight         = 1
	LoadBalanceTypeRandom     = "random"
	LoadBalanceTypeRoundRobin = "roundrobin"
)

// LoadBalance 负载均衡接口
type LoadBalance interface {
	Name() string
	Select(ctx context.Context, nodes []*registry.Node) (node *registry.Node, err error)
}

// NewLoadBalance 创建一个负载均衡器
// 参数
//   lbType: 负载均衡类型
// 返回值
//   LoadBalance: 负载均衡器
func NewLoadBalance(lbType string) LoadBalance {
	var lb LoadBalance
	switch lbType {
	case LoadBalanceTypeRandom:
		lb = NewRandomBalance()
	case LoadBalanceTypeRoundRobin:
		lb = NewRoundRobinBalance()
	default:
		lb = NewRandomBalance()
	}
	return lb
}
