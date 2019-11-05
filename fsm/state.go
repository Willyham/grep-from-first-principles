package fsm

// State is a single possible state in the machine.
type State struct {
	value       string
	isAccepting bool
}

// NewState creates a new, non-accepting state with the given value.
func NewState(value string) State {
	return State{
		value:       value,
		isAccepting: false,
	}
}

// NewAcceptingState creates a new accepting state with the given value.
func NewAcceptingState(value string) State {
	return State{
		value:       value,
		isAccepting: true,
	}
}

// MakeAccepting ensures that a state is accepting.
func (s State) MakeAccepting() State {
	return NewAcceptingState(s.Value())
}

// Accepting returns true if the state is accepting, false otherwise.
func (s State) Accepting() bool {
	return s.isAccepting
}

// Value of the state.
func (s State) Value() string {
	return s.value
}

// Equal checks equality based on the state value.
func (s State) Equal(other State) bool {
	return s.value == other.value
}
