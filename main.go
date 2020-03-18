package main

import (
	"fmt"

	graph "github.com/ram16230/compis1/Graph"
)

func delete(slice []string, el string) (a []string) {
	i := -1
	for j, s := range slice {
		if s == el {
			i = j
		}
	}
	a = append(slice[:i], slice[i+1:]...)
	return a
}

func main() {
	// Set Evaluator
	var ev Evaluator
	ev.operators = []string{"*", "+", ".", "|", "?"}
	ev.agrupation = []string{"()"}
	getList, alphabet := ev.separator("(a|b)*.a.b.b") // Get stack and alphabet
	ev.alphabet = alphabet                            // set alphabet
	fmt.Printf("GetList%v\n", getList)
	// Build Tree
	node := ev.getTree(getList)
	printTree(node)

	//AFN
	sigmaNotEpsilon := delete(ev.alphabet, "'")
	fmt.Printf("Alphabet%v\n", alphabet)
	fmt.Printf("sigmas %v\n", sigmaNotEpsilon)
	afnTree := InOrder(node, []string{"a", "b"})
	afn := afnTree.GetValue().(*graph.Automata)
	PrettyPrint(afn, "afn")

	//AFD
	afd := graph.NewAFDfromAFN(afn)
	PrettyPrint(afd, "afd")

	//AFDD
	// Set Again Evaluator
	ev = Evaluator{}
	ev.operators = []string{"*", "+", ".", "|", "?"}
	ev.agrupation = []string{"()"}
	getList, alphabet = ev.separator("(a|b)*.a.b.b") // Get stack and alphabet
	ev.alphabet = alphabet                           // set alphabet
	fmt.Printf("GetList%v\n", getList)
	// Build Tree Again
	n := ev.getTree(getList)
	printTree(n)
	// AFDD
	graph.IDTreeSet()(*n)
	afdd := graph.NewAFD(*n, []string{"a", "b"})
	PrettyPrint(afdd, "afdd")

	// Simulate
	t := "ab"
	fmt.Printf("afn: %v\n", afd.Simulate(t))
	fmt.Printf("afd: %v\n", afd.Simulate(t))
	fmt.Printf("afdd: %v\n", afdd.Simulate(t))
}
