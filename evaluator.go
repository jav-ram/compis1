package main // TODO: change main

import (
	"fmt"
	"strings"
)

type Evaluator struct {
	alphabet   []string
	operators  []string
	agrupation []string
}

type Token struct {
	name  string
	value string
}

// IsOpeningAgrupation gets if is an opening
func (ev *Evaluator) IsOpeningAgrupation(l byte) bool {
	return ev.IsAgrupation(string(l))
}

// IsClosingAgrupation gets if is an opening
func (ev *Evaluator) IsClosingAgrupation(l byte) bool {
	for _, op := range ev.agrupation {
		if l == op[1] {
			return true
		}
	}
	return false
}

// IsAgrupation tells if the lex contains a agrupation
func (ev *Evaluator) IsAgrupation(lex string) bool {
	for _, op := range ev.agrupation {
		if strings.Contains(lex, string(op[0])) {
			return true
		}
	}
	return false
}

// IsOperator tells if the lex is a operator or not
func (ev *Evaluator) IsOperator(lex string) bool {
	for _, op := range ev.operators {
		if op == lex {
			return true
		}
	}
	return false
}

// IsAlphabet tells if the lex is a operator or not
func (ev *Evaluator) IsAlphabet(lex string) bool {
	for _, w := range ev.alphabet {
		if w == lex {
			return true
		}
	}
	return false
}

// Init makes a function that can evalute
func (ev *Evaluator) Init() func(string) []interface{} {
	if ev.operators == nil {
		ev.operators = []string{"*", "+", "", "|"}
	}

	if ev.agrupation == nil {
		ev.agrupation = []string{"()"}
	}

	var separator func(input string) []interface{}
	separator = func(input string) (result []interface{}) {
		tmp := ""
		group := 0
		for i := 0; i < len(input); i++ {
			if ev.IsAgrupation(string(input[i])) {
				group++
				j := i + 1
				t := ""
				for j < len(input) {
					fmt.Printf("index: %d\n", j)
					if ev.IsOpeningAgrupation(input[j]) {
						group++
					} else if ev.IsClosingAgrupation(input[j]) {
						group--
					}
					if group <= 0 {
						break
					}
					t += string(input[j])
					fmt.Printf("letra: %s\n", string(input[j]))
					fmt.Printf("texto: %s\n", t)
					fmt.Printf("groups: %d\n", group)
					fmt.Println("--------------------------")
					j++
				}
				// delete substring
				i = j
				fmt.Printf("%v\n")
				result = append(result, separator(t))
			} else {
				tmp += string(input[i])
				if ev.IsOperator(tmp) || ev.IsAlphabet(tmp) {
					result = append(result, tmp)
					tmp = ""
				}
			}
		}
		return result
	}
	return separator
}

func main() {
	var ev Evaluator
	ev.alphabet = []string{"0", "1"}
	s := ev.Init()
	fmt.Printf("%v\n", s("(0|(10))0*(01)"))
}
