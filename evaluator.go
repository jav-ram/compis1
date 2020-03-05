package main // TODO: change main

import (
	"fmt"
	"strings"

	graph "github.com/ram16230/compis1/Graph"
	tree "github.com/ram16230/compis1/Tree"
)

type Evaluator struct {
	alphabet   []string
	operators  []string
	agrupation []string
}

type Token struct {
	name  string
	value string
}

func remove(s []interface{}, i int) []interface{} {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

// IsOpeningAgrupation gets if is an opening
func (ev *Evaluator) IsOpeningAgrupation(l byte) bool {
	return ev.IsAgrupation(string(l))
}

// IsClosingAgrupation gets if is an opening
func (ev *Evaluator) IsClosingAgrupation(l byte) bool {
	for _, op := range ev.agrupation {
		if l == op[1] {
			return true
		}
	}
	return false
}

// IsAgrupation tells if the lex contains a agrupation
func (ev *Evaluator) IsAgrupation(lex string) bool {
	for _, op := range ev.agrupation {
		if strings.Contains(lex, string(op[0])) {
			return true
		}
	}
	return false
}

// IsOperator tells if the lex is a operator or not
func (ev *Evaluator) IsOperator(lex string) bool {
	for _, op := range ev.operators {
		if op == lex {
			return true
		}
	}
	return false
}

// IsAlphabet tells if the lex is a operator or not
func (ev *Evaluator) IsAlphabet(lex string) bool {
	for _, w := range ev.alphabet {
		if w == lex {
			return true
		}
	}
	return false
}

func (ev *Evaluator) separator(input string) (result []interface{}) {
	tmp := ""
	group := 0
	for i := 0; i < len(input); i++ {
		if ev.IsOpeningAgrupation(input[i]) {
			group++
			j := i + 1
			t := ""
			for j < len(input) {
				fmt.Printf("index: %d\n", j)
				if ev.IsOpeningAgrupation(input[j]) {
					group++
				} else if ev.IsClosingAgrupation(input[j]) {
					group--
				}
				if group <= 0 {
					break
				}
				t += string(input[j])
				j++
			}
			// delete substring
			i = j
			result = append(result, ev.separator(t))
		} else {
			tmp += string(input[i])
			if ev.IsOperator(tmp) || ev.IsAlphabet(tmp) {
				result = append(result, tmp)
				tmp = ""
			}
		}
	}
	return result
}

func GetAutomata(node *tree.Node) *tree.Node {
	switch node.GetValue().(type) {
	case graph.Automata:
		{
			return node
		}
	default:
		{
			value := fmt.Sprintf("%v", node.GetValue())
			single := graph.SingleAFN(nil, value)
			node.SetValue(single)
			return node
		}
	}
}

func search(list []interface{}, s interface{}) int {
	for i, l := range list {
		if l == s {
			return i
		}
	}
	return -1
}

func (ev *Evaluator) getTree(input []interface{}) *tree.Node {

	for i := 0; i < len(input); i++ {
		switch input[i].(type) {
		case []interface{}:
			{
				v := input[i].([]interface{})
				input[i] = ev.getTree(v)
			}
		case *tree.Node:
			{
				break
			}
		case interface{}:
			{
				if ev.IsOperator(input[i].(string)) {
					c := input[i].(string)

					for search(input, "*") > 0 {
						n := tree.NewOpNode("*")
						idx := search(input, "*")
						l := ev.getTree([]interface{}{input[idx-1]})
						l.SetParent(n)
						n.AddLeftChild(l)
						input[idx-1] = n
						input = remove(input, idx)
					}

					for search(input, "+") > 0 {
						n := tree.NewOpNode("+")
						idx := search(input, "+")
						l := ev.getTree([]interface{}{input[idx-1]})
						l.SetParent(n)
						n.AddLeftChild(l)
						input[idx-1] = n
						input = remove(input, idx)
					}

					for search(input, "?") > 0 {
						n := tree.NewOpNode("?")
						idx := search(input, "?")
						l := ev.getTree([]interface{}{input[idx-1]})
						l.SetParent(n)
						n.AddLeftChild(l)
						input[idx-1] = n
						input = remove(input, idx)
					}

					for search(input, ".") > 0 {
						n := tree.NewOpNode(c)
						idx := search(input, ".")
						l := ev.getTree([]interface{}{input[idx-1]})
						l.SetParent(n)
						r := ev.getTree([]interface{}{input[idx+1]})
						r.SetParent(n)
						n.AddChilds(l, r)
						input[idx-1] = n
						input = remove(input, idx+1)
						input = remove(input, idx)
					}

					for search(input, "|") > 0 {
						n := tree.NewOpNode(c)
						idx := search(input, "|")
						l := ev.getTree([]interface{}{input[idx-1]})
						l.SetParent(n)
						r := ev.getTree([]interface{}{input[idx+1]})
						r.SetParent(n)
						n.AddChilds(l, r)
						input[idx-1] = n
						input = remove(input, idx+1)
						input = remove(input, idx)
					}

				} else {
					n := tree.NewLxNode(input[i].(string))
					input[i] = n
				}
			}
		}
	}
	return input[0].(*tree.Node)
}

func InOrder(node *tree.Node) *tree.Node {
	if node == nil {
		return nil
	}
	l := InOrder(node.Lchild)
	//r := InOrder(node.Rchild)
	if node.IsOperation() {
		switch node.GetValue() {
		case "*":
			{
				l = GetAutomata(l)
				newValue := graph.NewAFNKlean(nil, l.GetValue().(*graph.Automata))
				node.SetValue(newValue)
				return node
			}

		}
	}
	return node
}

func main() {
	/* var ev Evaluator
	ev.operators = []string{"*", "+", ".", "|", "?"}
	ev.alphabet = []string{"0", "1"}
	ev.agrupation = []string{"()"}
	getList := ev.separator("(0*)*")
	fmt.Printf("%v\n", getList)
	node := ev.getTree(getList)
	fmt.Printf("%v\n", node)
	afn := InOrder(node)
	fmt.Printf("%v\n", afn.GetValue())
	a := afn.GetValue().(*graph.Automata)
	r := a.Emulate("00")
	fmt.Printf("%v\n", r) */
	a := graph.SingleAFN([]string{"a"}, "a")
	aklean := graph.NewAFNKlean(nil, a)
	b := graph.SingleAFN(nil, "b")

	c := graph.NewAFNKOr(nil, aklean, b)
	op := graph.NewAFNKlean(nil, c)
	r := op.Emulate("i")
	fmt.Printf("%v\n", r)
	//op := graph.NewAFNKOr([]string{"1"}, a, b)
	//klean := graph.NewAFNConcat([]string{"a"}, op, a)

}
