package Graph

import (
	"fmt"
	"strconv"

	tree "github.com/ram16230/compis1/Tree"
)

// Nullable returns if the root contains epsilon on the language
func Nullable(root *tree.Node) bool {
	l := root.Lchild
	r := root.Rchild

	if l != nil && r != nil { // if it has two child
		return Nullable(l) || Nullable(r)
	} else if r == nil && l == nil { // if there is not a node
		return false
	} else if r != nil && l == nil { // if there is only one child
		return Nullable(l)
	}

	return false
}

// IDTree numera los valores
func IDTree(root *tree.Node, i int) tree.Node {
	v := root.GetValue().(string)

	node := tree.NewNode(v, false, false)

	if root.Lchild != nil {
		left := IDTree(root.Lchild, i)
		node.Lchild = &left
	}
	if root.Rchild != nil {
		right := IDTree(root.Rchild, i)
		node.Rchild = &right
	}

	if root.Lchild == nil && root.Rchild == nil {
		m := map[string][]string{}
		m["i"] = []string{strconv.Itoa(i)}
		m["v"] = []string{root.GetValue().(string)}
		t := fmt.Sprintf("%v", m)
		root.SetValue(t)
		i++
		fmt.Printf("%v\n", i)
	}

	return *node
}
