package main

import (
	"fmt"
	"io/ioutil"

	evaluator "github.com/ram16230/compis1/Evaluator"
	graph "github.com/ram16230/compis1/Graph"
	parser "github.com/ram16230/compis1/Parser"
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
	sigmaNotEpsilon := delete(ev.Alphabet, "`")
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
	// keywords
	tkns = append(tkns, token.NewKeywordTokenDescriptor("compiler", "COMPILER"))
	tkns = append(tkns, token.NewKeywordTokenDescriptor("characters", "CHARACTERS"))
	tkns = append(tkns, token.NewKeywordTokenDescriptor("keywords", "KEYWORDS"))
	tkns = append(tkns, token.NewKeywordTokenDescriptor("tokens", "TOKENS"))
	tkns = append(tkns, token.NewKeywordTokenDescriptor("productions", "PRODUCTIONS"))
	tkns = append(tkns, token.NewKeywordTokenDescriptor("chr", "CHR"))
	tkns = append(tkns, token.NewKeywordTokenDescriptor("except", "EXCEPT KEYWORDS"))
	// Tokens
	tkns = append(tkns, token.NewTokenDescriptor("ident", fmt.Sprintf("%v(%v|%v)*", letter, letter, digit))) // ident -> letter.(letter|digit)*
	tkns = append(tkns, token.NewTokenDescriptor("number", fmt.Sprintf("(%v)(%v*)", digit, digit)))          // number -> digit.(digit)*
	tkns = append(tkns, token.NewTokenDescriptor("equal", "="))                                              // equal -> =
	tkns = append(tkns, token.NewTokenDescriptor("finish", "."))                                             // finish -> .
	tkns = append(tkns, token.NewTokenDescriptor("complete", ".."))                                          // complete -> ..
	tkns = append(tkns, token.NewTokenDescriptor("group_start", "\\("))                                      // group_start -> (
	tkns = append(tkns, token.NewTokenDescriptor("group_end", "\\)"))                                        // group_end -> (
	tkns = append(tkns, token.NewTokenDescriptor("option_start", "["))                                       // option_start -> [
	tkns = append(tkns, token.NewTokenDescriptor("option_end", "]"))                                         // option_end -> ]
	tkns = append(tkns, token.NewTokenDescriptor("iteration_start", "{"))                                    // iteration_start -> {
	tkns = append(tkns, token.NewTokenDescriptor("iteration_end", "}"))                                      // iteration_edn -> }
	tkns = append(tkns, token.NewTokenDescriptor("quote", `"|'`))                                            // quote -> "
	tkns = append(tkns, token.NewTokenDescriptor("sum", "\\+"))                                              // sum -> +
	tkns = append(tkns, token.NewTokenDescriptor("subtract", "\\-"))                                         // subtract -> -
	tkns = append(tkns, token.NewTokenDescriptor("or", "|"))                                                 // subtract -> -
	scan := scanner.MakeAFNS(tkns)
	data, _ := ioutil.ReadFile("./test/DoubleAritmetica.ATG")

	tokens := scan.Simulate(string(data))

	parser.Parse(tokens)

}
