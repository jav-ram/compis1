package Scanner

import (
	"fmt"

	evaluator "github.com/ram16230/compis1/Evaluator"
	graph "github.com/ram16230/compis1/Graph"
	token "github.com/ram16230/compis1/Token"
)

type ScannerAF struct {
	automaton *graph.Automata
	name      string
	F         graph.Set
	Keywords  graph.Set
	sigma     []string
}

type ScannerAFCombined struct {
	automaton graph.Automata
	F         map[*graph.Automata][]string
	Keywords  map[*graph.Automata][]string
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

func MakeAFN(tkn token.TokenDescriptor) *ScannerAF {
	ev := evaluator.Evaluator{}
	ev.Operators = []string{"*", "+", "|", "?"}
	ev.Agrupation = []string{"()"}
	getList, alphabet := ev.Separator(tkn.Rgx) // Get stack and alphabet
	fmt.Printf("list: %v \n %v\n", getList, alphabet)
	ev.Alphabet = alphabet // set alphabet
	node := ev.GetTree(getList)

	// AFN
	sigmaNotEpsilon := delete(ev.Alphabet, "`")
	// fmt.Printf("sigmas %v\n", sigmaNotEpsilon)
	afnTree := evaluator.InOrder(node, sigmaNotEpsilon)
	afn := afnTree.GetValue().(*graph.Automata)
	// afd := graph.NewAFDfromAFN(afn) FIXME:

	keywords := *graph.NewSet()
	if tkn.IsKeyword {
		keywords = afn.F
	}

	// Return new Automaton
	scannerAFN := &ScannerAF{
		afn,
		tkn.ID,
		afn.F,
		keywords,
		ev.Alphabet,
	}
	return scannerAFN
}

// AddInParrallel
func AddInParrallel(auts ...*ScannerAF) *ScannerAFCombined {
	start := graph.NewAutomata(auts[0].sigma)
	r := &ScannerAFCombined{}
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
					t := map[string][]*graph.Automata{"`": []*graph.Automata{}} //TODO: epsilon
					t["`"] = append(t["`"], initial)                            //TODO: epsilon
					start.Trans[final] = t
				} else {
					t := start.Trans[final]          //TODO: epsilon
					t["`"] = append(t["`"], initial) //TODO: epsilon
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
	r.Keywords = map[*graph.Automata][]string{}
	// AddAllFinals and map them
	for _, aut := range auts {
		for final := range aut.automaton.F.GetList() {
			// ADD to table
			if r.F[final] == nil {
				t := []string{aut.name} //TODO: epsilon
				r.F[final] = t
				if len(aut.Keywords.GetItems()) > 0 {
					r.Keywords[final] = t
				}
			} else {
				t := r.F[final]
				t = append(t, aut.name) //TODO: epsilon
				r.F[final] = t
				if len(aut.Keywords.GetItems()) > 0 {
					r.Keywords[final] = t
				}
			}
			// ADD to automaton
			r.automaton.F.Add(final)
		}
	}

	return r
}

// MakeAFNS make an afn
func MakeAFNS(tkns []token.TokenDescriptor) *ScannerAFCombined {
	scannerAFs := []*ScannerAF{}
	for i := range tkns {
		scannerAFs = append(scannerAFs, MakeAFN(tkns[i]))
	}

	result := AddInParrallel(scannerAFs...)

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

func (scanner *ScannerAFCombined) Simulate(text string) []token.Token {
	aut := &scanner.automaton
	successStates := scanner.F
	keywordsStates := scanner.Keywords

	// init
	S := aut.Eclouser(&aut.Qo, graph.NewSet()) // current state
	var lexeme string                          // lexeme
	stack := []*graph.Set{}                    // Stack of acceptence states
	rememberIndex := []int{}                   // Stack of where in text was the acceptence state
	rememberToken := []string{}                //Stack of the lexeme of the acceptance state
	tokens := []token.Token{}

	// read
	i := -1
	c := text[0]

	// Main Scanner
	for len(text) > 0 {
		// while there is text
		for {
			if len(S.GetItems()) > 0 {
				// if State is not a state error
				i++
				c = text[i]
				// ----- move -----
				m := aut.Move(S, string(c))
				S = aut.Eclouser(m, graph.NewSet())
				// ----------------
				lexeme += string(text[i])
				g := contains(successStates, S)
				if g != nil {
					// if S is a goal state
					k := contains(keywordsStates, S)
					if k != nil {
						// if this state is of an keyword
						// it does not have more possible moves
						stack = []*graph.Set{}     // reset stack
						rememberIndex = []int{}    // reset index
						rememberToken = []string{} // reset lexeme
						// ----- We accept this as a token -----
						lex := text[:i+1]                                  // truncate lex of token
						text = text[i+1:]                                  // get rest of text
						tokens = append(tokens, token.NewToken(k[0], lex)) // Add to token list
						i = -1                                             // reset index of text
						S = aut.Eclouser(&aut.Qo, graph.NewSet())          // reset automaton
						if len(text) == 0 {
							// if there is no more break
							break
						}
						continue
						// ------------------------------------
					}

					// look if next char causes an error
					thereIsMore := false
					if len(text) > i+1 {
						m := aut.Move(S, string(text[i+1]))    // make a tmp next move FIXME: check if text has i + 1 elements
						tmp := aut.Eclouser(m, graph.NewSet()) // next state
						if len(tmp.GetItems()) > 0 {
							thereIsMore = true
						}
					}
					if thereIsMore {
						// it has more possible moves
						stack = append(stack, S)                      // stack state of acceptance
						rememberIndex = append(rememberIndex, i)      // stack index position of text for state of acceptance
						rememberToken = append(rememberToken, lexeme) // stack lexeme for the state of acceptance
					} else {
						// it does not have more possible moves
						stack = []*graph.Set{}     // reset stack
						rememberIndex = []int{}    // reset index
						rememberToken = []string{} // reset lexeme
						// ----- We accept this as a token -----
						lex := text[:i+1]                                  // truncate lex of token
						text = text[i+1:]                                  // get rest of text
						tokens = append(tokens, token.NewToken(g[0], lex)) // Add to token list
						i = -1                                             // reset index of text
						S = aut.Eclouser(&aut.Qo, graph.NewSet())          // reset automaton
						if len(text) == 0 {
							// if there is no more break
							break
						}
						// ------------------------------------
					}

				}
				/* fmt.Printf("S %v \n", S)
				fmt.Printf("tokens %v \n", tokens) */
			} else {
				// if there is an error
				if len(stack) > 0 {
					// there is a state in stack
					S = stack[len(stack)-1]                         // pop last state of acceptance
					stack = []*graph.Set{}                          // reset stack
					i = rememberIndex[len(rememberIndex)-1]         // pop last index of acceptance
					rememberIndex = []int{}                         // reset index
					newToken := rememberToken[len(rememberToken)-1] // pop last index of acceptance
					rememberToken = []string{}                      // reset lexeme
					// make token
					// ----- We accept this as a token -----
					lex := text[:i+1]                                      // truncate lex of token
					text = text[i+1:]                                      // get rest of text
					tokens = append(tokens, token.NewToken(newToken, lex)) // Add to token list
					i = -1                                                 // reset index of text
					S = aut.Eclouser(&aut.Qo, graph.NewSet())              // reset automaton
					if len(text) == 0 {
						// if there is no more break
						break
					}
					// ------------------------------------
				} else {
					// syntax error
					stack = []*graph.Set{}     // reset stack
					rememberIndex = []int{}    // reset index
					rememberToken = []string{} // reset lexeme
					// ----- We accept this as a token -----
					text = text[i+1:]                         // get rest of text
					i = -1                                    // reset index of text
					S = aut.Eclouser(&aut.Qo, graph.NewSet()) // reset automaton
					if len(text) == 0 {
						// if there is no more break
						break
					}
					// ------------------------------------
				}
			}
		}
	}

	fmt.Printf("tokens: \n%v\n", tokens)

	return tokens

}
