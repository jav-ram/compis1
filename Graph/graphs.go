package Graph

import "fmt"

// Automata is the struct of a automata machine
type Automata struct {
	Q     Set
	Sigma []string
	Qo    Set
	Trans map[*Automata]map[string][]*Automata
	F     Set
}

// State automata
type State Automata

func NewAutomata() *Automata {
	a := &Automata{}
	a.Trans = map[*Automata]map[string][]*Automata{}
	a.Qo = *NewSetFrom(a)
	a.Q = *NewSetFrom(a)
	a.F = *NewSetFrom(a)
	return a
}

// AFN automata
type AFN Automata

// AFD automata
type AFD Automata

func SetList(qs []*Automata, count int) []*Automata {
	for i := 0; i < count; i++ {
		qs = append(qs, NewAutomata())
	}
	return qs
}

//AddBySide
func (first *Automata) AddBySide(second *Automata) *Automata {
	trans := map[*Automata]map[string][]*Automata{}

	for t1, v := range first.Trans {
		for c1, j := range v {
			if trans[t1] == nil {
				t := map[string][]*Automata{c1: j}
				trans[t1] = t
			} else {
				t := trans[t1]
				t[c1] = append(t[c1], j...)
				trans[t1] = t
			}
		}
	}

	for t2, v := range second.Trans {
		for c2, j := range v {
			if trans[t2] == nil {
				t := map[string][]*Automata{c2: j}
				trans[t2] = t
			} else {
				t := trans[t2]
				t[c2] = append(t[c2], j...)
				trans[t2] = t
			}
		}
	}

	for tf := range first.F.list {
		for to := range second.Qo.list {
			if trans[tf] == nil {
				t := map[string][]*Automata{"": []*Automata{to}}
				trans[tf] = t
			} else {
				t := trans[tf]
				t[""] = append(t[""], to)
				trans[tf] = t
			}
		}
	}

	return &Automata{
		*Union(&first.Q, &second.Q),
		first.Sigma,
		first.Qo,
		trans,
		second.F,
	}
}

//SingleAFN
func SingleAFN(sigma []string, va string) *Automata {
	var qs []*Automata
	qs = SetList(qs, 2)

	q := NewSetFrom(qs...)
	qo := NewSetFrom(qs[0])
	qf := NewSetFrom(qs[1])

	trans := map[*Automata]map[string][]*Automata{}
	t := map[string][]*Automata{va: []*Automata{qs[1]}}
	trans[qs[0]] = t

	a := NewAutomata()
	a.Q = *q
	a.Sigma = sigma
	a.Qo = *qo
	a.Trans = trans
	a.F = *qf
	return a
}

//NewAFNKlean get Klean AFN
func NewAFNKlean(sigma []string, aut *Automata) *Automata {
	s := NewAutomata()
	l := NewAutomata()
	r := NewAutomata()

	r = s.AddBySide(aut)
	r = r.AddBySide(l)

	// outside
	for to := range s.Qo.list {
		for tf := range l.F.list {
			fmt.Printf("aaaaaaaaaaaaa: %v\n", r.Trans[to])
			if r.Trans[to] == nil {
				t := map[string][]*Automata{"": []*Automata{tf}}
				r.Trans[to] = t
			} else {
				t := r.Trans[to]
				t[""] = append(t[""], tf)
				r.Trans[to] = t
			}
			fmt.Printf("aaaaaaaaaaaaa: %v\n", r.Trans[to])
		}
	}

	// inner
	for tf := range aut.F.list {
		for to := range aut.Qo.list {
			t := r.Trans[tf]
			t[""] = append(t[""], to)
			r.Trans[tf] = t
		}
	}

	fmt.Printf("trans %v\n", r.Trans)

	return r
}

//NewAFNKConcat get Concat AFN
func NewAFNKConcat(sigma []string, first *Automata, second *Automata) *Automata {
	r := NewAutomata()
	r = first.AddBySide(second)
	return r
}
