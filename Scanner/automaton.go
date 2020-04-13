package Scanner

import (
	graph "github.com/ram16230/compis1/Graph",
	evaluator "github.com/ram16230/compis1/Evaluator",
)

type ScannerAF struct {
	automaton *graph.Automata,
	name string,
	F []*graph.Automata
}

type ScannerAFCombined struct {
	automaton graph.Automata,
	F map[*graph.Automata][]string{}
}

func MakeAFN(
	ev evaluator.Ev,
	name string,
	sigma []string,
) ScannerAF {
	node := ev.GetTree(getList)

	// AFN
	sigmaNotEpsilon := delete(ev.Alphabet, "'")
	fmt.Printf("Alphabet%v\n", alphabet)
	// fmt.Printf("sigmas %v\n", sigmaNotEpsilon)
	afnTree := evaluator.InOrder(node, sigmaNotEpsilon)
	afn := afnTree.GetValue().(*graph.Automata)
	
	// Return new Automaton
	scannerAFN := ScannerAF{
		afn,
		name,
		afn.Qo.list,
	}
	return scannerAFN
}

// AddInParrallel
func AddInParrallel(sigma []string, auts ...*ScannerAF) *ScannerAFCombined {
	start := graph.NewAutomata(sigma)
	r := *ScannerAFScannerAFCombined{}

	// Make All conections from start to new automaton
	for _, aut := range auts {
		start.Q.Adds(aut.automaton.Q)
		// for every automata
		for final : range start.F.list {
			// for every final state of start
			for initial range aut.Qo.list {
				// for every intial start of aut

				// make a transition from final to initial with sigma
				t := map[string][]*Automata{"'": []*Automata{}} //TODO: epsilon
				t["'"] = append(t["'"], initial)  //TODO: epsilon
				start.trans[final] = t
			}
		}
	}

	r.automaton = start

	// AddAllFinals and map them
	for _, aut := range auts {
		for final : range aut.automaton.F.list {
			if r.F[final] == nil {
				t := []string{}
				t = append(t, aut.name)
				r.F[final] = t
			} else {
				r.F[final] = append(r.F[final], aut.name)
			}
		}
	}

	return r
}

func MakeAFNS(
	evs []evaluator.Evaluator,
	names []string,
	sigma []string,
) *ScannerAFCombined {
	scannerAFs := []ScannerAF{}
	for i := range evs {
		scannerAFs = append(scannerAFs, MakeAFN(evs[i], names[i], sigma))
	}

	return AddInParrallel(sigma, ...evs)
}

func (afn ScannerAF) AFD()
