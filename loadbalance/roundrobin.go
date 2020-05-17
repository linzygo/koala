package loadbalance

import (
	"context"
	"koala/errno"
	"koala/logger"
	"koala/registry"
)

// NewRoundRobinBalance 创建一个轮询负载均衡器
func NewRoundRobinBalance() LoadBalance {
	return &RoundRobinBalance{}
}

// RoundRobinBalance 负载均衡，轮询算法
type RoundRobinBalance struct {
	curWeight int // 当前权重值
}

// Name 实现接口LoadBalance
func (rr *RoundRobinBalance) Name() string {
	return "roundrobin"
}

// Select 实现接口LoadBalance
func (rr *RoundRobinBalance) Select(ctx context.Context, nodes []*registry.Node) (node *registry.Node, err error) {
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
		if node.Weight == 0 {
			v.Weight = DefaultNodeWeight // node是指针类型，这里修改会导致内存改变，下次遍历得到的是这个值
		}
		totalWeight += v.Weight
	}

	// 这个轮询其实是有问题的，因为得到的服务节点都是一样的，这时候所有客户端都是先访问第一个节点，有瞬间压垮的可能
	curWight := rr.curWeight
	for _, v := range newNodes {
		curWight -= node.Weight
		if curWight < 0 {
			node = v
			rr.curWeight = (rr.curWeight + 1) % totalWeight
			return
		}
	}

	err = errno.InvalidNode
	logger.Error(ctx, "找不到有效节点")
	return
}
