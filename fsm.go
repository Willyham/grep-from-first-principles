package fsm

/*
Sigma  is the input alphabet (a finite, non-empty set of symbols). -> ASCII
S is a finite, non-empty set of states.
s0 is an initial state, an element of S.
Delta is the state-transition function: Delta: Sigma x S -> S
F is the set of final states, a (possibly empty) subset of {\displaystyle S} S.
*/

type State string

type transition struct {
	name      string
	source    State
	nextState State
}

type StateMachine struct {
	currentState State
	finalStates  []State
	transitions  []transition
}

func New(initial State, final []State, transitions []transition) (*StateMachine, error) {
	// TODO: Ensure set of valid transitions
	return &StateMachine{
		currentState: initial,
		finalStates:  final,
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
		m.currentState = transition.nextState
	}
	// TODO: Check for success if no final states given.
	// After all input, check if the state we finished in is in the set of final states.
	for _, s := range m.finalStates {
		if m.currentState == s {
			return true
		}
	}
	return false
}

func (m *StateMachine) findTransition(next string) *transition {
	for _, t := range m.transitions {
		if t.source == m.currentState && t.name == next {
			return &t
		}
	}
	return nil
}
