package main

import (
	"fmt"
	caller "github.com/ram16230/compis1/Caller"
)

func a(i int) int {
	fmt.Printf("%v\n", caller.GetCallerFunctionName())
	return i
}

func add(i int) int {
	return i
}

func addOne(v *int) {
	current := 0
	_ = current

	// blablablbalbalblabla
	// ...
	// ...

	*v = add(a(*v))
	*v = 1111
}

func main() {
	var a int = 0
	addOne(&a)
	fmt.Printf("%v\n", a)
}
