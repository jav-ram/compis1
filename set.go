package main // TODO: change main

// Set is a struct that defines a set
type Set struct {
	list map[int]bool
}

// NewSet returns a new set
func NewSet() *Set {
	s := &Set{}
	s.list = make(map[int]bool)
	return s
}

// Has saids if i is in the set
func (set *Set) Has(i int) bool {
	_, ok := set.list[i]
	return ok
}

// Add adds i into the set
func (set *Set) Add(i int) {
	set.list[i] = true
}

// Remove removes i inside set
func (set *Set) Remove(i int) {
	delete(set.list, i)
}

// Adds adds multiples int to set
func (set *Set) Adds(I ...int) {
	for _, i := range I {
		set.Add(i)
	}
}

// Removes removes multiples int of set
func (set *Set) Removes(I ...int) {
	for _, i := range I {
		set.Remove(i)
	}
}

// Union gets two sets and returns a new set with the values of both sets
func Union(A *Set, B *Set) *Set {
	s := NewSet()
	for a, v := range A.list {
		if v {
			s.Add(a)
		}
	}

	for b, v := range A.list {
		if v {
			s.Add(b)
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
