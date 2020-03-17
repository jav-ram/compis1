package Graph

import (
	"fmt"
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

func contains(D []*Automata, S *Automata) bool {
	for _, d := range D {
		if reflect.DeepEqual(d.ids, S.ids) {
			return true
		}
	}
	return false
}

func NewAFD(iroot tree.Node, sigma []string) *Automata {
	newAut := NewAutomata(sigma)
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

	// set initial states
	so := FirstPos(*root)
	qo := NewAutomata(sigma)
	qo.ids = so
	newAut.Qo = qo.Q

	D := []*Automata{}
	D = append(D, qo)

	// set F
	F := ""
	for n, node := range nmaps {
		switch c := node.GetValue().(type) {
		case string:
			{
				if c == "#" {
					F = n
				}
			}
		case map[string][]string:
			{
				if c["v"][0] == "#" {
					F = n
				}
			}
		}
	}

	for UnMark(D) {
		S := GetUnMark(D)
		S.mark = true

		for _, a := range sigma {
			U := []string{}
			u := NewAutomata(sigma)

			for _, p := range S.ids {
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

			u.ids = U

			if !contains(D, u) && len(U) != 0 {
				D = append(D, u)
				// is final state ?
				for _, h := range u.ids {
					if h == F {
						newAut.F.Add(u)
					}
				}
			}

			tmp := NewAutomata(sigma)
			tmp.ids = U

			var eq *Automata
			for _, t := range D {
				fmt.Printf("t %v u %v\n", t.ids, U)
				if reflect.DeepEqual(t.ids, U) {
					fmt.Println("Son iguales")
					eq = t
				}
			}

			if eq != nil {
				u = eq
			}

			// add transition from S to U with a
			fmt.Printf("%v[%v]=> %v \n", S.ids, a, u.ids)
			if newAut.Trans[S] == nil {
				y := map[string][]*Automata{a: []*Automata{u}}
				newAut.Trans[S] = y
			} else {
				y := newAut.Trans[S]
				y[a] = append(y[a], u)
				newAut.Trans[S] = y
			}

		}

	}
	newAut.Q = *NewSetFrom(D...)
	return newAut
}
