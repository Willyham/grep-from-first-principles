package fsm

import (
	"github.com/awalterschulze/gographviz"
)

/*
Sigma  is the input alphabet (a finite, non-empty set of symbols). -> ASCII
S is a finite, non-empty set of states.
s0 is an initial state, an element of S.
Delta is the state-transition function: Delta: Sigma x S -> S
F is the set of final states, a (possibly empty) subset of {\displaystyle S} S.
*/

type Transition struct {
	Event     string
	Source    State
	NextState State
}

type StateMachine struct {
	currentState State
	transitions  []Transition
}

func New(initial State, transitions []Transition) (*StateMachine, error) {
	// TODO: Ensure set of valid transitions
	return &StateMachine{
		currentState: initial,
		transitions:  transitions,
	}, nil
}

func (m *StateMachine) Run(events []string) bool {
	for _, event := range events {
		transition := m.findTransition(event)
		if transition == nil {
			// if no transition exists between current state and next character, there is no match.
			return false
		}
		m.currentState = transition.NextState
	}

	// After all input, check if the state we finished in is in the set of final states.
	return m.currentState.Accepting()
}

func (m *StateMachine) findTransition(next string) *Transition {
	for _, t := range m.transitions {
		if t.Source.Equal(m.currentState) && t.Event == next {
			return &t
		}
	}
	return nil
}

func (m *StateMachine) ToGraphViz() string {
	graph := gographviz.NewGraph()
	graph.Name = "FSM"
	graph.Directed = true

	uniqueStates := map[State]bool{}
	for _, transition := range m.transitions {
		seenSource := uniqueStates[transition.Source]
		if !seenSource {
			uniqueStates[transition.Source] = true
			attrs := map[string]string{}
			if transition.Source.Accepting() {
				attrs[string(gographviz.Shape)] = "doublecircle"
			}
			graph.AddNode(graph.Name, stringOrEpsilon(transition.Source.Value()), attrs)
		}
		seenDest := uniqueStates[transition.NextState]
		if !seenDest {
			uniqueStates[transition.NextState] = true
			attrs := map[string]string{}
			if transition.NextState.Accepting() {
				attrs[string(gographviz.Shape)] = "doublecircle"
			}
			graph.AddNode(graph.Name, stringOrEpsilon(transition.NextState.Value()), attrs)
		}
		graph.AddEdge(stringOrEpsilon(transition.Source.Value()), stringOrEpsilon(transition.NextState.Value()), true, map[string]string{
			"label": stringOrEpsilon(transition.Event),
		})
	}
	return graph.String()
}

func stringOrEpsilon(in string) string {
	if in == "" {
		return "<&#949;>"
	}
	return in
}
