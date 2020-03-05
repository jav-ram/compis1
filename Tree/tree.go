package Tree // TODO: change main

type Node struct {
	parent      *Node
	Lchild      *Node
	Rchild      *Node
	value       interface{}
	isOperation bool
	isLexema    bool
	Read        bool
}

func NewNode(value interface{}, isOperation bool, isLexema bool) *Node {
	n := &Node{}
	n.value = value
	n.isOperation = isOperation
	n.isLexema = isLexema
	return n
}

func (node *Node) IsOperation() bool {
	return node.isOperation
}

func (node *Node) IsLexema() bool {
	return node.isLexema
}

func (child *Node) GetParent() *Node {
	return child.parent
}

func (child *Node) SetParent(parent *Node) {
	child.parent = parent
}

func (root *Node) GetValue() interface{} {
	return root.value
}

func (root *Node) SetValue(value interface{}) {
	root.value = value
}

func NewOpNode(value interface{}) *Node {
	return NewNode(value, true, false)
}

func NewLxNode(value interface{}) *Node {
	return NewNode(value, false, true)
}

func (root *Node) AddLeftChild(node *Node) {
	root.Lchild = node
	node.parent = root
}

func (root *Node) AddRightChild(node *Node) {
	root.Rchild = node
	node.parent = root
}

func (root *Node) AddChilds(lNode *Node, rNode *Node) {
	root.AddLeftChild(lNode)
	root.AddRightChild(rNode)
}

func (node *Node) GetRoot() *Node {
	for node.parent != nil {
		node = node.parent
	}
	return node
}
