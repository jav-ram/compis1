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

func SimulateAll(t string, rg string) {
	// Set Evaluator
	var ev Evaluator
	ev.operators = []string{"*", "+", ".", "|", "?"}
	ev.agrupation = []string{"()"}
	getList, alphabet := ev.separator(rg) // Get stack and alphabet
	ev.alphabet = alphabet                // set alphabet
	fmt.Printf("GetList%v\n", getList)
	// Build Tree
	node := ev.getTree(getList)
	printTree(node)

	//AFN
	sigmaNotEpsilon := delete(ev.alphabet, "'")
	fmt.Printf("Alphabet%v\n", alphabet)
	fmt.Printf("sigmas %v\n", sigmaNotEpsilon)
	afnTree := InOrder(node, sigmaNotEpsilon)
	afn := afnTree.GetValue().(*graph.Automata)
	PrettyPrint(afn, "afn", rg)

	//AFD
	fmt.Printf("afn sigma%v\n", afn.Sigma)
	afn.Sigma = sigmaNotEpsilon
	afd := graph.NewAFDfromAFN(afn)
	PrettyPrint(afd, "afd", rg)

	//AFDD
	// Set Again Evaluator
	ev = Evaluator{}
	ev.operators = []string{"*", "+", ".", "|", "?"}
	ev.agrupation = []string{"()"}
	getList, alphabet = ev.separator(rg) // Get stack and alphabet
	ev.alphabet = alphabet               // set alphabet
	fmt.Printf("GetList%v\n", getList)
	// Build Tree Again
	n := ev.getTree(getList)
	printTree(n)
	// AFDD
	graph.IDTreeSet()(*n)
	afdd := graph.NewAFD(*n, sigmaNotEpsilon)
	PrettyPrint(afdd, "afdd", rg)

	// Simulate
	fmt.Printf("-afn: %v\n", afd.Simulate(t))
	fmt.Printf("-afd: %v\n", afd.Simulate(t))
	fmt.Printf("afdd: %v\n", afdd.Simulate(t))
}

func main() {
	regexs := []string{
		"(a*|b*).c",
		"(b|b)*.a.b.b.(a|b)*",
		"(a|')*.b.(a.a*).c",
		"(a|b)*.a.(a|b).(a|b)",
		"b*.a.(b|')",
		"a.(b)*.a.(b)*",
		"(('|0).1*)*",
		"(0|1)*.0.(0|1).(0|1)",
		"(0.0)*.(1.1)*",
		"(0|1).1*.(0|1)",
	}
	goods := []string{
		"aaaaaaaaaaaaaaac",
		"babba",
		"baac",
		"bbbbaab",
		"a",
		"abbbbbbbab",
		"011011",
		"1111111001",
		"000011",
		"011111",
	}

	for i := range regexs {
		SimulateAll(goods[i], regexs[i])
	}

	//SimulateAll("ab", "b*.a.b")
}
