package registry

// Service 服务抽象
type Service struct {
	Name  string  `json:"name"`
	Nodes []*Node `json:"nodes"`
}

// Node 服务节点抽象
type Node struct {
	ID     string `json:"id"`
	IP     string `json:"ip"`
	Port   int    `json:"port"`
	Weight int    `json:"weight"`
}

// NodeInfo 单个服务节点的信息, Name为服务名称
type NodeInfo struct {
	ServiceName string `json:"name"`
	Node        `json:"node"`
}
