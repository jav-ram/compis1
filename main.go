package main

import (
	"fmt"

	evaluator "github.com/ram16230/compis1/Evaluator"
	graph "github.com/ram16230/compis1/Graph"
	scanner "github.com/ram16230/compis1/Scanner"
)

// TODO: move to utils
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

func SimulateAll(t string, rg string) {
	// Set Evaluator
	var ev evaluator.Evaluator
	ev.Operators = []string{"*", "+", ".", "|", "?"}
	ev.Agrupation = []string{"()"}
	getList, alphabet := ev.Separator(rg) // Get stack and alphabet
	ev.Alphabet = alphabet                // set alphabet
	fmt.Printf("GetList%v\n", getList)
	// Build Tree
	node := ev.GetTree(getList)
	//printTree(node)

	//AFN
	sigmaNotEpsilon := delete(ev.Alphabet, "'")
	fmt.Printf("Alphabet%v\n", alphabet)
	// fmt.Printf("sigmas %v\n", sigmaNotEpsilon)
	afnTree := evaluator.InOrder(node, sigmaNotEpsilon)
	afn := afnTree.GetValue().(*graph.Automata)
	// PrettyPrint(afn, "afn", rg)

	//AFD
	fmt.Printf("afn sigma%v\n", afn.Sigma)
	afn.Sigma = sigmaNotEpsilon
	afd := graph.NewAFDfromAFN(afn)
	//PrettyPrint(afd, "afd", rg)

	//AFDD
	// Set Again Evaluator
	ev = evaluator.Evaluator{}
	ev.Operators = []string{"*", "+", ".", "|", "?"}
	ev.Agrupation = []string{"()"}
	getList, alphabet = ev.Separator(rg) // Get stack and alphabet
	ev.Alphabet = alphabet               // set alphabet
	fmt.Printf("GetList%v\n", getList)
	// Build Tree Again
	n := ev.GetTree(getList)
	// AFDD
	graph.IDTreeSet()(*n)
	afdd := graph.NewAFD(*n, sigmaNotEpsilon)
	evaluator.PrettyPrint(afdd, "afdd", rg)

	// Simulate
	fmt.Printf("%v => %v\n", t, rg)
	fmt.Printf("-afn: %v\n", afd.Simulate(t))
	fmt.Printf("-afd: %v\n", afd.Simulate(t))
	fmt.Printf("afdd: %v\n", afdd.Simulate(t))
}

func main() {
	SimulateAll("abb", "(b|b)*.a.b.b.(a|')*")
	rgs := []string{
		"a.b*.a",
		"b.a*.b",
	}
	names := []string{
		"A",
		"B",
	}
	scanner.MakeAFNS(rgs, names).Simulate("abbbbbaaabaaaab")
}
