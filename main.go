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

	// SimulateAll("abb", "(b|b)*.a.b.b.(a|')*")
	rgs := []string{
		fmt.Sprintf("%v(%v|%v)*", letter, letter, digit), // ident -> letter.(letter|digit)*
		fmt.Sprintf("(%v)(%v*)", digit, digit),           // number -> digit.(digit)*
		"=",                                              // equal -> =
		".",                                              // finish -> .
		fmt.Sprintf("\\%v", "("),                         // group -> (
		fmt.Sprintf("\\%v", ")"),                         // group -> (
		fmt.Sprintf("\\%v", "["),                         // group -> [
		fmt.Sprintf("\\%v", "]"),                         // group -> ]
		fmt.Sprintf("\\%v", "{"),                         // group -> {
		fmt.Sprintf("\\%v", "}"),                         // group -> }
	}
	names := []string{
		"ident",
		"number",
		"equal",
		"finish",
		"group start",
		"group end",
		"option start",
		"option end",
		"iteration start",
		"iteration end",
	}
	/* ev := evaluator.Evaluator{}
	ev.Operators = []string{"*", "+", "|", "?"}
	ev.Agrupation = []string{"()"}
	getList, _ := ev.Separator(rgs[2]) // Get stack and alphabet
	fmt.Printf("List %v\n", getList) */
	scanner.MakeAFNS(rgs, names).Simulate("(abs123=123.")
	/* ev := evaluator.Evaluator{}
	ev.Operators = []string{"*", "+", "|", "?"}
	ev.Agrupation = []string{"()"}
	getList, alphabet := ev.Separator(rgs[2]) // Get stack and alphabet
	fmt.Printf("list: %v \n%v\n", getList, alphabet)
	ev.Alphabet = alphabet // set alphabet
	node := ev.GetTree(getList)

	evaluator.PrintTree(node)

	// AFN
	sigmaNotEpsilon := delete(ev.Alphabet, "'")
	// fmt.Printf("sigmas %v\n", sigmaNotEpsilon)
	afnTree := evaluator.InOrder(node, sigmaNotEpsilon)
	afn := afnTree.GetValue().(*graph.Automata)

	evaluator.PrettyPrint(afn, "afn", "oyea")

	fmt.Printf("%v\n", afn) */

}
