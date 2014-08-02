package cmsgparser

// NodeType defines the type of Node
type NodeType int

const (
	cantbe NodeType = iota
	TextNode
	LinkNode
	MoveActionNode
	RootNode
)

// String returns string representation of this NodeType
func (t NodeType) String() string {
	switch t {
	case TextNode:
		return "TextNode"
	case LinkNode:
		return "LinkNode"
	case MoveActionNode:
		return "MoveActionNode"
	case RootNode:
		return "RootNode"
	}
	return "Not defined type error!"
}

type Node struct {
	Type     NodeType
	Value    string
	children []*Node
}

func (n *Node) String() string {
	return n.Value
}

func (n *Node) Children() []*Node {
	return n.children
}

func (n *Node) AddChild(child *Node) {
	n.children = append(n.children, child)
}
