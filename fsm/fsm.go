package fsm

import (
	"github.com/awalterschulze/gographviz"
)

// Transition describes a possible way to move between states.
type Transition struct {
	Event     string
	Source    State
	NextState State
}

// StateMachine allows us to play events and compute the resulting state.
type StateMachine struct {
	initialState State
	currentState State
	transitions  []Transition
}

// New creates a StateMachine from the given transitions and initial state.
// TODO: Ensure that transitions are valid by:
// - Checking that all events exist within a given alphabet
// - Ensuring that the list of transitions is deterministic (i.e. no duplicate
//   [source,event] tuples)
func New(initial State, transitions []Transition) *StateMachine {
	return &StateMachine{
		initialState: initial,
		currentState: initial,
		transitions:  transitions,
	}
}

// Run events through the state machine and determine if the resulting state is accepting.
func (m *StateMachine) Run(events []string) bool {
	for _, event := range events {
		transition := m.findTransition(event)
		if transition == nil {
			// If no transition exists between current state and next character, we're done.
			// Test if the state we're currently in is accepting.
			return m.currentState.Accepting()
		}
		m.currentState = transition.NextState
	}

	// After all input, check if the state we finished in is in the set of final states.
	return m.currentState.Accepting()
}

// Reset the state machine to its initial state.
func (m *StateMachine) Reset() {
	m.currentState = m.initialState
}

func (m *StateMachine) findTransition(event string) *Transition {
	for _, t := range m.transitions {
		if t.Source.Equal(m.currentState) && t.Event == event {
			return &t
		}
	}
	return nil
}

// ToGraphViz generates a graphviz representation of the state machine.
func (m *StateMachine) ToGraphViz() string {
	graph := gographviz.NewGraph()
	graph.Name = "FSM"
	graph.Directed = true
	graph.Attrs[gographviz.RankDir] = "LR"

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
	attrs[gographviz.Shape] = "circle"
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
