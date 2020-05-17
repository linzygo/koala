package loadbalance

import (
	"context"
	"koala/errno"
	"koala/logger"
	"koala/registry"
	"math/rand"
)

// NewRandomBalance 创建一个随机负载均衡器
func NewRandomBalance() LoadBalance {
	return &RandomBalance{}
}

// RandomBalance 负载均衡，权重随机算法
type RandomBalance struct {
}

// Name 实现接口LoadBalance
func (r *RandomBalance) Name() string {
	return "random"
}

// Select 实现接口LoadBalance
func (r *RandomBalance) Select(ctx context.Context, nodes []*registry.Node) (node *registry.Node, err error) {
	if len(nodes) == 0 {
		logger.Error(ctx, "传入的节点队列为空")
		err = errno.EmptyNode
		return
	}

	newNodes := filterNodes(ctx, nodes)
	if len(newNodes) == 0 {
		logger.Error(ctx, "所有节点都失败")
		return
	}

	defer func() {
		if node != nil {
			setSelectedNode(ctx, node)
		}
	}()

	totalWeight := 0
	for _, v := range newNodes {
		if v.Weight == 0 {
			v.Weight = DefaultNodeWeight // node是指针类型，这里修改会导致内存改变，下次遍历得到的是这个值
		}
		totalWeight += v.Weight
	}
	curWeight := rand.Intn(totalWeight)
	for _, v := range newNodes {
		curWeight -= v.Weight
		if curWeight < 0 {
			node = v
			return
		}
	}

	err = errno.InvalidNode
	logger.Error(ctx, "找不到有效节点")
	return
}
