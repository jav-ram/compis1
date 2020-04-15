package Scanner

import (
	"fmt"

	evaluator "github.com/ram16230/compis1/Evaluator"
	graph "github.com/ram16230/compis1/Graph"
)

type ScannerAF struct {
	automaton *graph.Automata
	name      string
	F         graph.Set
	sigma     []string
}

type ScannerAFCombined struct {
	automaton graph.Automata
	F         map[*graph.Automata][]string
}

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

func MakeAFN(
	rg string,
	name string,
) *ScannerAF {
	ev := evaluator.Evaluator{}
	ev.Operators = []string{"*", "+", ".", "|", "?"}
	ev.Agrupation = []string{"()"}
	getList, alphabet := ev.Separator(rg) // Get stack and alphabet
	ev.Alphabet = alphabet                // set alphabet
	node := ev.GetTree(getList)

	// AFN
	sigmaNotEpsilon := delete(ev.Alphabet, "'")
	// fmt.Printf("sigmas %v\n", sigmaNotEpsilon)
	afnTree := evaluator.InOrder(node, sigmaNotEpsilon)
	afn := afnTree.GetValue().(*graph.Automata)

	// Return new Automaton
	scannerAFN := &ScannerAF{
		afn,
		name,
		afn.Qo,
		ev.Alphabet,
	}
	return scannerAFN
}

// AddInParrallel
func AddInParrallel(auts ...*ScannerAF) *ScannerAFCombined {
	start := graph.NewAutomata(auts[0].sigma)
	r := &ScannerAFCombined{}

	fmt.Printf("finals of start %v\n", start.F.GetList())
	fmt.Printf("iniciales de As %v\n", auts[0].automaton.Qo.GetList())
	fmt.Printf("iniciales de Bs %v\n", auts[1].automaton.Qo.GetList())
	// Make All conections from start to new automaton
	for _, aut := range auts {
		start.Q.Adds(aut.automaton.Q.GetItems()...)

		// for every automata
		for final := range start.F.GetList() {
			// for every final state of start
			for initial := range aut.automaton.Qo.GetList() {
				// for every intial start of aut

				// make a transition from final to initial with sigma
				if start.Trans[final] == nil {
					t := map[string][]*graph.Automata{"'": []*graph.Automata{}} //TODO: epsilon
					t["'"] = append(t["'"], initial)                            //TODO: epsilon
					start.Trans[final] = t
				} else {
					t := start.Trans[final]          //TODO: epsilon
					t["'"] = append(t["'"], initial) //TODO: epsilon
					start.Trans[final] = t
				}
			}
		}

		for from, withTo := range aut.automaton.Trans {
			// for every transition on aut
			if start.Trans[from] == nil {
				start.Trans[from] = withTo
			} else {
				for with, to := range withTo {
					if withTo[with] == nil {
						t := to
						start.Trans[from][with] = t
					} else {
						t := start.Trans[from][with]
						t = append(t, to...)
						start.Trans[from][with] = t
					}
				}
			}
		}
	}

	fmt.Printf("Transitions %v\n", start.Trans)

	r.automaton = *start
	r.F = map[*graph.Automata][]string{}
	// AddAllFinals and map them
	for _, aut := range auts {
		for final := range aut.automaton.F.GetList() {
			// ADD to table
			if r.F[final] == nil {
				t := []string{aut.name} //TODO: epsilon
				r.F[final] = t
			} else {
				t := r.F[final]
				t = append(t, aut.name) //TODO: epsilon
				r.F[final] = t
			}
			// ADD to automaton
			r.automaton.F.Add(final)
		}
	}

	return r
}

func MakeAFNS(
	rgs []string,
	names []string,
) *ScannerAFCombined {
	scannerAFs := []*ScannerAF{}
	for i := range rgs {
		scannerAFs = append(scannerAFs, MakeAFN(rgs[i], names[i]))
	}

	result := AddInParrallel(scannerAFs...)
	evaluator.PrettyPrint(&result.automaton, "afdd", "ohyea")

	return result
}

func contains(tabSuccessStates map[*graph.Automata][]string, set *graph.Set) []string {
	for _, s := range set.GetItems() {
		if tabSuccessStates[s] != nil {
			return tabSuccessStates[s]
		}
	}
	return nil
}

func (scanner *ScannerAFCombined) Simulate(text string) {
	aut := &scanner.automaton
	successStates := scanner.F

	// init
	S := aut.Eclouser(&aut.Qo, graph.NewSet()) // current state
	var lexeme string                          // lexeme
	dollar := graph.NewAutomata([]string{})    // $
	dollar.ToggleMark()                        // set as $
	stack := []*graph.Automata{dollar}         // FIXME: get global "$"
	tokens := []string{}

	// read
	i := -1
	c := text[0]

	// Main Scanner
	for len(text) > 0 {
		fmt.Printf("Outer for %v\n", S)
		for len(S.GetItems()) > 0 { // While State is not a state of error
			i++
			c = text[i]
			stack = append(stack, S.GetItems()...)
			// move
			m := aut.Move(S, string(c))
			S = aut.Eclouser(m, graph.NewSet())
			fmt.Printf("FFFFFF %v %v\n", text, i)
			fmt.Printf("G %v %v\n", text, i)
			lexeme += string(text[i])
			g := contains(successStates, S)
			if g != nil {
				fmt.Println("Is sucess")
				// if S is a goal state
				stack = []*graph.Automata{} // clear
				lex := text[:i+1]           // truncate
				text = text[i+1:]
				tokens = append(tokens, fmt.Sprintf("<%v, %v>", g, lex))
				i = -1
				S = aut.Eclouser(&aut.Qo, graph.NewSet())
				if len(text) == 0 {
					break
				}
			}
			fmt.Printf("tokens %v \n", tokens)
		}
		// Roll back Loop TODO: make it work
		for contains(successStates, S) == nil && len(stack) > 0 {
			x := &graph.Automata{}
			fmt.Printf("stack %v \n", len(stack))
			x, stack = stack[len(stack)-1], stack[:len(stack)-1] // pop
			S = graph.NewSetFrom(x)
			lexeme = lexeme[:len(lexeme)-1] // delete last char
		}
	}

	fmt.Printf("tokens: \n%v\n", tokens)

}
