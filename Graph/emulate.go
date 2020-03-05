package Graph

import "fmt"

func (aut *Automata) eclouser(state *Set, visited *Set) *Set {
	r := NewSet()

	fmt.Printf("k: %v\n", state.list)
	for k := range state.list {
		r.Adds(k)
		for c, a := range aut.Trans[k] {
			if c == "" {
				b := []*Automata{}
				for _, t := range a {
					if !visited.Has(t) {
						visited.Add(t)
						b = append(b, t)
					}
				}
				s := NewSetFrom(b...)
				r = Union(r, aut.eclouser(s, visited))
			}
		}
	}

	return r
}

func (aut *Automata) move(state *Set, t string) *Set {
	r := NewSet()

	for k := range state.list {
		for c, a := range aut.Trans[k] {
			if c == t {
				r.Adds(a...)
			}
		}
	}

	return r
}

func (aut *Automata) Emulate(text string) *Set {
	S := aut.eclouser(&aut.Qo, NewSet())
	for _, c := range text {
		m := aut.move(S, string(c))
		fmt.Println("------------------")
		fmt.Printf("aut q: %v\n", aut.Q)
		fmt.Printf("aut trans: %v\n", aut.Trans)
		fmt.Printf("s: %v\n", S)
		fmt.Printf("c: %v\n", string(c))
		fmt.Printf("m: %v\n", m)
		S = aut.eclouser(m, NewSet())
		fmt.Printf("s: %v\n", S)
		fmt.Printf("F: %v\n", aut.F)
		fmt.Println("------------------")
	}
	return Intersection(S, &aut.F)
}

func main() {

}
