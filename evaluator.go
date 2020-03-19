package main // TODO: change main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	gographvix "github.com/awalterschulze/gographviz"
	gotree "github.com/disiqueira/gotree"
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
	return append(s[:i], s[i+1:]...)
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

func (ev *Evaluator) separator(input string) (result []interface{}, alphabet []string) {
	group := 0
	for i := 0; i < len(input); i++ {
		if ev.IsOpeningAgrupation(input[i]) {
			group++
			j := i + 1
			t := ""
			for j < len(input) {
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
			sep, a := ev.separator(t)
			alphabet = append(alphabet, a...)
			result = append(result, sep)
		} else {
			c := string(input[i])
			if ev.IsOperator(c) {
				result = append(result, c)
			} else {
				result = append(result, c)
				alphabet = append(alphabet, c)
			}
		}
	}
	alphabet = append(alphabet, "'")
	return result, unique(alphabet)
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
						n := tree.NewOpNode(".")
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

func InOrder(node *tree.Node, sigma []string) *tree.Node {
	if node == nil {
		return nil
	}
	l := InOrder(node.Lchild, sigma)
	r := InOrder(node.Rchild, sigma)
	if node.IsOperation() {
		switch node.GetValue() {
		// Single Operators
		case "*":
			{
				lAut := l.GetValue().(*graph.Automata)
				newValue := graph.NewAFNKlean(sigma, lAut)
				node.SetValue(newValue)
				node.Lchild = nil
				return node
			}
		case "+":
			{
				lAut := l.GetValue().(*graph.Automata)
				newValue := graph.NewAFNSum(sigma, lAut)
				node.SetValue(newValue)
				node.Lchild = nil
				return node
			}
		case "?":
			{
				lAut := l.GetValue().(*graph.Automata)
				newValue := graph.NewAFNQuestion(sigma, lAut)
				node.SetValue(newValue)
				return node
			}
		case "|":
			{
				lAut := l.GetValue().(*graph.Automata)
				rAut := r.GetValue().(*graph.Automata)
				newValue := graph.NewAFNKOr(sigma, lAut, rAut)
				node.SetValue(newValue)
				return node
			}
		case ".":
			{
				lAut := l.GetValue().(*graph.Automata)
				rAut := r.GetValue().(*graph.Automata)
				newValue := graph.NewAFNConcat(sigma, lAut, rAut)
				node.SetValue(newValue)
				return node
			}
		}
	} else {
		return GetAutomata(node)
	}
	return node
}

func unique(strSlice []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range strSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

func PrettyPrint(aut *graph.Automata, dir string, name string) {
	f, err := os.Create("python/graph.txt")
	if err != nil {
		fmt.Println(err)
		f.Close()
		return
	}
	value := aut.Trans

	G := gographvix.NewGraph()
	G.Attrs.Add("rankdir", "LR")
	G.SetName("G")
	G.SetDir(true)

	tmp := map[*graph.Automata]int{}
	qo := []int{}
	qf := []int{}
	i := 0
	for k := range value {
		_, ok := tmp[k]
		if !ok {
			tmp[k] = i
			if aut.Qo.Has(k) {
				a := map[string]string{}
				a["shape"] = "square"
				G.AddNode("G", strconv.Itoa(i), a)
				qo = append(qo, i)
			} else if aut.F.Has(k) {
				a := map[string]string{}
				a["shape"] = "doublecircle"
				G.AddNode("G", strconv.Itoa(i), a)
				qf = append(qf, i)
			} else {
				G.AddNode("G", strconv.Itoa(i), nil)
			}
			i++
		}
		for c, is := range value[k] {
			if c == "'" { //TODO: epsilon
				c = "ɛ"
			}
			for _, j := range is {
				_, ok := tmp[j]
				if !ok {
					tmp[j] = i
					if aut.Qo.Has(j) {
						a := map[string]string{}
						a["shape"] = "square"
						G.AddNode("G", strconv.Itoa(i), a)
						qo = append(qo, i)
					} else if aut.F.Has(j) {
						a := map[string]string{}
						a["shape"] = "doublecircle"
						G.AddNode("G", strconv.Itoa(i), a)
						qf = append(qf, i)
					} else {
						G.AddNode("G", strconv.Itoa(i), nil)
					}
					i++
				}
				a := map[string]string{}
				a["label"] = c
				G.AddEdge(strconv.Itoa(tmp[k]), strconv.Itoa(tmp[j]), true, a)
				// fmt.Printf("x: [%v][%s] => [%v]\n", tmp[k], c, tmp[j])
				t := fmt.Sprintf("%v,%v,%v", tmp[k], c, tmp[j])
				fmt.Fprintln(f, t)
			}
		}
	}
	fs, err := os.Create(fmt.Sprintf("graphs/%v/%v.dot", dir, name))
	if err != nil {
		fmt.Println(err)
		fs.Close()
		return
	}
	fmt.Fprintln(fs, G.String())
	//fmt.Printf("FFFFFF: %v\n", aut.F)
}

func getOtherNode(root *tree.Node) gotree.Tree {

	t := fmt.Sprintf("%v", root.GetValue())

	node := gotree.New(t)

	if root.Lchild != nil {
		left := getOtherNode(root.Lchild)
		node.AddTree(left)
	}
	if root.Rchild != nil {
		right := getOtherNode(root.Rchild)
		node.AddTree(right)
	}

	return node
}

func printTree(root *tree.Node) {
	gtree := getOtherNode(root)
	fmt.Printf("------------------------\n")
	fmt.Println(gtree.Print())
	fmt.Printf("------------------------\n")
}
