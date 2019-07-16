package fsm

import (
	"github.com/awalterschulze/gographviz"
)

var DefaultAlphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

type Transition struct {
	Event     string
	Source    State
	NextState State
}

type StateMachine struct {
	initialState State
	currentState State
	transitions  []Transition
	alphabet     []string
}

func New(initial State, transitions []Transition, alphabet []string) (*StateMachine, error) {
	return &StateMachine{
		initialState: initial,
		currentState: initial,
		transitions:  transitions,
		alphabet:     alphabet,
	}, nil
}

func (m *StateMachine) Run(events []string) bool {
	for _, event := range events {
		transition := m.findTransition(event)
		if transition == nil {
			// if no transition exists between current state and next character, we're done.
			// Test if the state we're currently in is accepting.
			return m.currentState.Accepting()
		}
		m.currentState = transition.NextState
	}

	// After all input, check if the state we finished in is in the set of final states.
	return m.currentState.Accepting()
}

func (m *StateMachine) Reset() {
	m.currentState = m.initialState
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
		nodes := transitionToNodes(transition)
		seenSource := uniqueStates[transition.Source]
		if !seenSource {
			graph.Nodes.Add(&nodes[0])
		}
		seenNext := uniqueStates[transition.NextState]
		if !seenNext {
			graph.Nodes.Add(&nodes[1])
		}
		graph.AddEdge(stringOrEpsilon(transition.Source.Value()), stringOrEpsilon(transition.NextState.Value()), true, map[string]string{
			"label": stringOrEpsilon(transition.Event),
		})
	}
	return graph.String()
}

func transitionToNodes(transition Transition) []gographviz.Node {
	return []gographviz.Node{
		stateToNode(transition.Source),
		stateToNode(transition.NextState),
	}
}

func stateToNode(state State) gographviz.Node {
	attrs, _ := gographviz.NewAttrs(nil) // Cannot error if empty map
	if state.Accepting() {
		attrs[gographviz.Shape] = "doublecircle"
	}
	return gographviz.Node{
		Name:  stringOrEpsilon(state.Value()),
		Attrs: attrs,
	}
}

func stringOrEpsilon(in string) string {
	if in == "" {
		return "<&#949;>"
	}
	return in
}
