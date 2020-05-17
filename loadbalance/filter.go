package loadbalance

import (
	"context"
	"fmt"

	"koala/registry"
)

type selectedNodes struct {
	selectedNodeMap map[string]bool
}

type selectedNodesKey struct {
}

// WithSelectedNodes context保存选择过的节点
func WithSelectedNodes(ctx context.Context) context.Context {
	selNodes := &selectedNodes{
		selectedNodeMap: make(map[string]bool),
	}
	return context.WithValue(ctx, selectedNodesKey{}, selNodes)
}

func getSelectedNodes(ctx context.Context) *selectedNodes {
	selNodes, ok := ctx.Value(selectedNodesKey{}).(*selectedNodes)
	if !ok {
		selNodes = &selectedNodes{}
	}
	return selNodes
}

func setSelectedNode(ctx context.Context, node *registry.Node) {
	selNodes := getSelectedNodes(ctx)
	selNode := fmt.Sprintf("%s:%d", node.IP, node.Port)
	selNodes.selectedNodeMap[selNode] = true
}

// filterNodes 过滤出未被选择过的节点
// 参数
//   ctx: 带选择过的节点
//   nodes: 所有节点
// 返回值
//   []*registry.Node: 未用过的节点
func filterNodes(ctx context.Context, nodes []*registry.Node) []*registry.Node {
	var retNodes []*registry.Node
	selNodes := getSelectedNodes(ctx)
	for _, node := range nodes {
		k := fmt.Sprintf("%s:%d", node.IP, node.Port)
		if _, ok := selNodes.selectedNodeMap[k]; !ok {
			retNodes = append(retNodes, node)
		}
	}
	return retNodes
}
