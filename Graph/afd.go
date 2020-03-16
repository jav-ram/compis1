package Graph

import (
	"bufio"
	"fmt"
	"os"
	"reflect"
	"strconv"

	tree "github.com/ram16230/compis1/Tree"
)

// Nullable returns if the root contains epsilon on the language
func Nullable(root *tree.Node) bool {

	switch c := root.GetValue().(type) {
	case string:
		{
			if c == "#" {
				// epsilon
				return true
			} else if c == "|" {
				// c1 or c2
				return Nullable(root.Lchild) || Nullable(root.Rchild)
			} else if c == "." {
				// c1 and c2
				return Nullable(root.Lchild) && Nullable(root.Rchild)
			} else if c == "*" {
				// klean
				return true
			}
		}
	default:
		{
			// leaf
			return false
		}
	}
	return false
}

// IDTreeSet numera los valores
func IDTreeSet() func(tree.Node) (tree.Node, map[string]tree.Node) {
	i := 0
	nmaps := map[string]tree.Node{}

	var IDTree func(tree.Node) (tree.Node, map[string]tree.Node)
	IDTree = func(root tree.Node) (tree.Node, map[string]tree.Node) {

		var node tree.Node

		switch v := root.GetValue().(type) {
		case string:
			{
				node = *tree.NewNode(v, false, false)
			}
		case map[string][]string:
			{
				node = *tree.NewNode(v["v"], false, false)
			}
		}

		if root.Lchild != nil {
			l := root.Lchild
			if l.Lchild == nil && l.Rchild == nil {
				i++
				m := map[string][]string{}

				switch c := l.GetValue().(type) {
				case string:
					{
						m["v"] = []string{c}
					}
				case map[string][]string:
					{
						m["v"] = c["v"]
					}
				}

				m["i"] = []string{strconv.Itoa(i)}
				nmaps[strconv.Itoa(i)] = *l
				l.SetValue(m)
			}
			left, lmap := IDTree(*root.Lchild)
			for k, v := range lmap {
				nmaps[k] = v
			}
			node.Lchild = &left
		}
		if root.Rchild != nil {
			r := root.Rchild
			if r.Lchild == nil && r.Rchild == nil {
				i++
				m := map[string][]string{}

				switch c := r.GetValue().(type) {
				case string:
					{
						m["v"] = []string{c}
					}
				case map[string][]string:
					{
						m["v"] = c["v"]
					}
				}

				m["i"] = []string{strconv.Itoa(i)}
				nmaps[strconv.Itoa(i)] = *r
				r.SetValue(m)
			}
			right, rmap := IDTree(*root.Rchild)
			for k, v := range rmap {
				nmaps[k] = v
			}
			node.Rchild = &right
		}

		return node, nmaps
	}
	return IDTree
}

// FirstPos gets the list of ids of the first pos
func FirstPos(root tree.Node) []string {
	switch c := root.GetValue().(type) {
	case string:
		{
			if c == "'" {
				// epsilon
				return []string{}
			} else if c == "|" {
				// c1 or c2
				c1 := root.Lchild
				c2 := root.Rchild
				fs := []string{}
				fs = append(fs, FirstPos(*c1)...)
				fs = append(fs, FirstPos(*c2)...)
				// Union c1 c2
				return fs
			} else if c == "." {

				if Nullable(root.Lchild) {
					c1 := root.Lchild
					c2 := root.Rchild
					fs := []string{}
					fs = append(fs, FirstPos(*c1)...)
					fs = append(fs, FirstPos(*c2)...)
					// Union c1 c2
					return fs
				}

				c1 := root.Lchild
				return FirstPos(*c1)

			} else if c == "*" {
				// klean
				c1 := root.Lchild
				return FirstPos(*c1)
			}
		}
	case map[string][]string:
		{
			// leaf
			return c["i"]
		}
	}
	// leaf
	m := root.GetValue().(map[string][]string)
	return m["i"]
}

// LastPos gets the list of ids of the last pos
func LastPos(root tree.Node) []string {

	switch c := root.GetValue().(type) {
	case string:
		{
			if c == "'" {
				// epsilon
				return []string{}
			} else if c == "|" {
				// c1 or c2
				c1 := root.Lchild
				c2 := root.Rchild
				ls := []string{}
				ls = append(ls, LastPos(*c1)...)
				ls = append(ls, LastPos(*c2)...)
				// Union c1 c2
				return ls
			} else if c == "." {

				if Nullable(root.Rchild) {
					c1 := root.Lchild
					c2 := root.Rchild
					ls := []string{}
					ls = append(ls, LastPos(*c1)...)
					ls = append(ls, LastPos(*c2)...)
					// Union c1 c2
					return ls
				}

				c2 := root.Rchild
				return LastPos(*c2)

			} else if c == "*" {
				// klean
				c1 := root.Lchild
				return LastPos(*c1)
			}
		}
	case map[string][]string:
		{
			// leaf
			return c["i"]
		}
	}
	// leaf
	m := root.GetValue().(map[string][]string)
	return m["i"]
}

func SetFollowPos(root tree.Node) map[string][]string {
	if &root == nil {
		return nil
	}

	fmap := map[string][]string{}

	var v string
	switch c := root.GetValue().(type) {
	case string:
		{
			v = c
		}
	case map[string][]string:
		{
			v = c["v"][0]
		}
	}

	if root.Lchild != nil {
		m := SetFollowPos(*root.Lchild)
		for k, v := range m {
			_, ok := fmap[k]
			if !ok {
				fmap[k] = v
			} else {
				fmap[k] = append(fmap[k], v...)
			}
		}
	}

	if root.Rchild != nil {
		m := SetFollowPos(*root.Rchild)
		for k, v := range m {
			_, ok := fmap[k]
			if !ok {
				fmap[k] = v
			} else {
				fmap[k] = append(fmap[k], v...)
			}
		}
	}

	if v == "." {
		c1 := *root.Lchild
		c2 := *root.Rchild

		for _, i := range LastPos(c1) {
			fmap[i] = append(fmap[i], FirstPos(c2)...)
		}
	} else if v == "*" {
		for _, i := range LastPos(root) {
			fmap[i] = append(fmap[i], FirstPos(root)...)
		}
	}

	return fmap
}

func Intersect(A, B []string) []string {
	r := []string{}
	for _, a := range A {
		for _, b := range B {
			if a == b {
				r = append(r, a)
			}
		}
	}
	return r
}

func unmark(D [][]string, marks [][]string) []string {
	if len(marks) == 0 {
		return D[0]
	}
	for _, s := range D {
		if !contains(marks, s) {
			return s
		}
	}
	return nil
}

func contains(D [][]string, S []string) bool {
	for _, d := range D {
		if reflect.DeepEqual(d, S) {
			return true
		}
	}
	return false
}

func NewAFD(iroot tree.Node, sigma []string) *Automata {
	newAut := NewAutomata(sigma)
	newAut.Qo = *NewSet()
	newAut.Q = *NewSet()
	newAut.F = *NewSet()
	// init add concat #
	root := tree.NewOpNode(".")
	finish := tree.NewLxNode("#")
	root.Lchild = &iroot
	root.Rchild = finish
	_, nmaps := IDTreeSet()(*root)
	fmap := SetFollowPos(*root)
	fmt.Println("*****************************")
	fmt.Printf("%v\n", fmap)
	fmt.Println("*****************************")
	D := [][]string{}
	mark := [][]string{}

	fmt.Printf("%v\n", fmap)

	D = append(D, FirstPos(*root))

	// set initial states
	so := unmark(D, mark)
	qo := NewAutomata(sigma)
	qo.ids = so
	newAut.Qo.Add(qo)
	newAut.Q.Add(qo)

	for unmark(D, mark) != nil { // while unmark exists on list
		fmt.Printf("D: %v\nQ: %v\n", D, newAut.Q.list)
		fmt.Printf("Trans: %v", newAut.Trans)
		reader := bufio.NewReader(os.Stdin)
		reader.ReadString('\n')
		S := unmark(D, mark)
		s := NewAutomata(sigma)
		s.ids = S
		mark = append(mark, S)

		for _, a := range sigma {
			U := []string{}

			for _, p := range S {
				n := nmaps[p]
				if n.GetValue() == "#" {
					continue
				}
				nv := n.GetValue().(map[string][]string)

				for _, v := range nv["v"] {
					if v == a {
						U = append(U, fmap[p]...)
					}
				}
			}

			if !contains(D, U) && len(U) != 0 {
				D = append(D, U)
			}

			tmp := NewAutomata(sigma)
			tmp.ids = U

			var eq *Automata
			var u *Automata
			fmt.Println("**********************")
			for t := range newAut.Q.list {
				fmt.Printf("IDS%v ", t.ids)
			}
			fmt.Println("")
			fmt.Printf("%v\n", U)
			for t := range newAut.Q.list {
				if reflect.DeepEqual(t.ids, U) {
					eq = t
				}
			}

			if eq != nil {
				u = eq
				fmt.Printf("u: %v\n", u)
			} else {
				fmt.Println("no exists")
				u = NewAutomata(sigma)
				u.ids = U
				u.Q.Adds(u)
				fmt.Printf("u: %v\n", u)
			}
			fmt.Println("**********************")

			// add transition from S to U with a
			fmt.Printf("%v[%v]=> %v \n", s.ids, a, u.ids)
			if newAut.Trans[s] == nil {
				y := map[string][]*Automata{a: []*Automata{u}}
				newAut.Trans[s] = y
			} else {
				y := newAut.Trans[s]
				y[a] = append(y[a], u)
				newAut.Trans[s] = y
			}

		}

	}
	fmt.Printf("AAAAAAAAAAAAAAAAAAAAAa\n%v\n", newAut)
	return newAut
}
