package metricnode

import "strings"

type NodeContainer struct {
	Nodes []*MetricNode
}

type MetricNode struct {
	Value      string
	ChildNodes []*MetricNode
}

func NewMetricNode(value string) *MetricNode {
	return parseNode(value)
}

func NewNodeContainer() *NodeContainer {
	return &NodeContainer{}
}

func parseNode(value string) *MetricNode {
	splitValue := strings.SplitN(value, ".", 2)

	if len(splitValue) == 1 {
		return &MetricNode{
			Value:      splitValue[0],
			ChildNodes: []*MetricNode{},
		}
	}

	return &MetricNode{
		Value: splitValue[0],
		ChildNodes: []*MetricNode{
			parseNode(splitValue[1]),
		},
	}
}

func (n *NodeContainer) Store(node *MetricNode) {
	n.Nodes = rangeOverNodes(n.Nodes, node)
}

func (n *NodeContainer) MaxChildNodes() int {
	max := 1
	for _, node := range n.Nodes {
		subTotal := totalSubNodes(node, max)
		if subTotal > max {
			max = subTotal
		}
	}
	return max
}

func totalSubNodes(node *MetricNode, max int) int {
	max++
	if len(node.ChildNodes) > 0 {
		return MaxNodes(node.ChildNodes)
	}

	return max
}

func MaxNodes(nodes []*MetricNode) int {
	max := 1
	for _, node := range nodes {
		subTotal := totalSubNodes(node, max)
		if subTotal > max {
			max = subTotal
		}
	}
	return max
}

func rangeOverNodes(nodes []*MetricNode, node *MetricNode) []*MetricNode {
	nodeExists := false
	for _, n := range nodes {
		if n.Value == node.Value {
			nodeExists = true
			if len(n.ChildNodes) > 0 {
				n.ChildNodes = rangeOverNodes(n.ChildNodes, node.ChildNodes[0])
			} else {
				n.ChildNodes = append(n.ChildNodes, node)
				break
			}
		}
	}
	if !nodeExists {
		nodes = append(nodes, node)
	}

	return nodes
}
