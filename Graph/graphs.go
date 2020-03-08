package Graph

import (
	"fmt"
	"os"
)

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

func PrintAutomata(trans string) {
	f, err := os.Create("python/graph.txt")
	if err != nil {
		fmt.Println(err)
		f.Close()
		return
	}
	t := fmt.Sprintf("%v", trans)

	fmt.Fprintln(f, t)

}

// MergeTrans Merge two transitions
func MergeTrans(m1, m2 map[*Automata]map[string][]*Automata) map[*Automata]map[string][]*Automata {
	r := map[*Automata]map[string][]*Automata{}

	for k, m := range m1 {
		for mc, ma := range m {
			if r[k] == nil {
				t := map[string][]*Automata{mc: ma}
				r[k] = t
			} else {
				t := r[k]
				t[mc] = append(t[mc], ma...)
				r[k] = t
			}
		}
	}

	for k, m := range m2 {
		for mc, ma := range m {
			if r[k] == nil {
				t := map[string][]*Automata{mc: ma}
				r[k] = t
			} else {
				t := r[k]
				t[mc] = append(t[mc], ma...)
				r[k] = t
			}
		}
	}

	return r
}

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
			if r.Trans[to] == nil {
				t := map[string][]*Automata{"": []*Automata{tf}}
				r.Trans[to] = t
			} else {
				t := r.Trans[to]
				t[""] = append(t[""], tf)
				r.Trans[to] = t
			}
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

	return r
}

//NewAFNConcat get Concat AFN
func NewAFNConcat(sigma []string, first *Automata, second *Automata) *Automata {
	r := NewAutomata()
	r = first.AddBySide(second)
	return r
}

//NewAFNKOr get Concat AFN
func NewAFNKOr(sigma []string, a *Automata, b *Automata) *Automata {
	r := NewAutomata()
	s := NewAutomata()
	f := NewAutomata()

	r.Q = *NewSetFrom(s, f, a, b)
	r.Qo = *NewSetFrom(s)
	r.F = *NewSetFrom(f)
	r.Sigma = sigma

	// primeras dos trancisiones
	tmp1 := NewAutomata()
	tmp1 = s.AddBySide(a)
	tmp2 := NewAutomata()
	tmp2 = s.AddBySide(b)
	r.Trans = MergeTrans(r.Trans, a.Trans)
	r.Trans = MergeTrans(r.Trans, b.Trans)

	t := map[string][]*Automata{"": []*Automata{}}
	t[""] = append(t[""], tmp1.Trans[s][""]...)
	t[""] = append(t[""], tmp2.Trans[s][""]...)
	r.Trans[s] = t

	// segundas trancisión
	for kf := range a.F.list {
		for ko := range f.Qo.list {
			if r.Trans[kf] == nil {
				t := map[string][]*Automata{"": []*Automata{ko}}
				r.Trans[kf] = t
			} else {
				t := r.Trans[kf]
				t[""] = append(t[""], ko)
				r.Trans[kf] = t
			}
		}
	}

	for kf := range b.F.list {
		for ko := range f.Qo.list {
			if r.Trans[kf] == nil {
				t := map[string][]*Automata{"": []*Automata{ko}}
				r.Trans[kf] = t
			} else {
				t := r.Trans[kf]
				t[""] = append(t[""], ko)
				r.Trans[kf] = t
			}
		}
	}

	return r
}

// NewAFNSum un automata de +
func NewAFNSum(sigma []string, aut *Automata) *Automata {
	f := NewAFNKlean(sigma, aut)
	r := NewAFNConcat(sigma, aut, f)
	return r
}

// NewAFNQuestion un automata de ?
func NewAFNQuestion(sigma []string, aut *Automata) *Automata {
	f := SingleAFN(sigma, "")
	r := NewAFNKOr(sigma, aut, f)
	return r
}
