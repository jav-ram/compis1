package main

import (
	"fmt"

	evaluator "github.com/ram16230/compis1/Evaluator"
	graph "github.com/ram16230/compis1/Graph"
	scanner "github.com/ram16230/compis1/Scanner"
	token "github.com/ram16230/compis1/Token"
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

func MakeLetter() string {
	set := ""
	for i := 65; i <= 90; i++ {
		set += string(rune(i)) + "|"
	}
	for i := 97; i < 122; i++ {
		set += string(rune(i)) + "|"
	}
	set += string(rune(122))
	return "(" + set + ")"
}

func MakeDigit() string {
	set := ""
	for i := 0; i < 9; i++ {
		set += fmt.Sprintf("%v", i) + "|"
	}
	set += fmt.Sprintf("%v", 9)
	return "(" + set + ")"
}

func main() {
	letter := MakeLetter()
	digit := MakeDigit()

	var tkns []token.TokenDescriptor
	tkns = append(tkns, token.NewTokenDescriptor("ident", fmt.Sprintf("%v(%v|%v)*", letter, letter, digit))) // ident -> letter.(letter|digit)*
	tkns = append(tkns, token.NewTokenDescriptor("number", fmt.Sprintf("(%v)(%v*)", digit, digit)))          // number -> digit.(digit)*
	tkns = append(tkns, token.NewTokenDescriptor("equal", "="))                                              // equal -> =
	tkns = append(tkns, token.NewTokenDescriptor("finish", "."))                                             // finish -> .
	tkns = append(tkns, token.NewTokenDescriptor("group_start", fmt.Sprintf("\\%v", "(")))                   // group_start -> (
	tkns = append(tkns, token.NewTokenDescriptor("group_end", fmt.Sprintf("\\%v", ")")))                     // group_end -> (
	tkns = append(tkns, token.NewTokenDescriptor("option_start", fmt.Sprintf("\\%v", "[")))                  // option_start -> [
	tkns = append(tkns, token.NewTokenDescriptor("option_end", fmt.Sprintf("\\%v", "]")))                    // option_end -> ]
	tkns = append(tkns, token.NewTokenDescriptor("iteration_start", fmt.Sprintf("\\%v", "{")))               // iteration_start -> {
	tkns = append(tkns, token.NewTokenDescriptor("iteration_end", fmt.Sprintf("\\%v", "}")))                 // iteration_edn -> }

	scanner.MakeAFNS(tkns).Simulate("(abs123=123.")

}
