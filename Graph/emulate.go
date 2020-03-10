package Graph

func (aut *Automata) eclouser(state *Set, visited *Set) *Set {
	r := NewSet()

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
		S = aut.eclouser(m, NewSet())
	}
	return Intersection(S, &aut.F)
}
