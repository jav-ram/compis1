package Graph

func (aut *Automata) Eclouser(state *Set, visited *Set) *Set {
	r := NewSet()

	for k := range state.list {
		r.Adds(k)
		for c, a := range aut.Trans[k] {
			if c == "'" { // TODO: epsilon
				b := []*Automata{}
				for _, t := range a {
					if !visited.Has(t) {
						visited.Add(t)
						b = append(b, t)
					}
				}
				s := NewSetFrom(b...)
				r = Union(r, aut.Eclouser(s, visited))
			}
		}
	}

	return r
}

func (aut *Automata) Move(state *Set, t string) *Set {
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

func (aut *Automata) Simulate(text string) bool {
	S := aut.Eclouser(&aut.Qo, NewSet())
	for _, c := range text {
		m := aut.Move(S, string(c))
		S = aut.Eclouser(m, NewSet())
	}
	g := Intersection(S, &aut.F)
	if len(g.list) > 0 {
		return true
	}
	return false
}
