package Graph // TODO: change main

import (
	"fmt"
)

// Set is a struct that defines a set
type Set struct {
	list map[*Automata]bool
}

// NewSet returns a new set
func NewSet() *Set {
	s := &Set{}
	s.list = make(map[*Automata]bool)
	return s
}

func NewSetFrom(I ...*Automata) *Set {
	s := NewSet()
	for _, i := range I {
		s.Add(i)
	}
	return s
}

// Has saids if i is in the set
func (set *Set) Has(i *Automata) bool {
	fmt.Printf("iii: %v\n", i)
	_, ok := set.list[i]
	return ok
}

// Add adds i into the set
func (set *Set) Add(b *Automata) {
	set.list[b] = true
}

// Remove removes i inside set
func (set *Set) Remove(i *Automata) {
	delete(set.list, i)
}

// Adds adds multiples int to set
func (set *Set) Adds(I ...*Automata) {
	for _, i := range I {
		set.Add(i)
	}
}

// Removes removes multiples int of set
func (set *Set) Removes(I ...*Automata) {
	for _, i := range I {
		set.Remove(i)
	}
}

// Union gets two sets and returns a new set with the values of both sets
func Union(A *Set, B *Set) *Set {
	s := NewSet()

	for k, v := range A.list {
		if v {
			s.Add(k)
		}
	}

	for k, v := range B.list {
		if v {
			s.Add(k)
		}
	}

	return s
}

// Intersection gets two sets and returns a new set the keys that are inside both sets
func Intersection(A *Set, B *Set) *Set {
	s := NewSet()
	for k, v := range A.list {
		if v && B.Has(k) {
			s.Add(k)
		}
	}
	return s
}
