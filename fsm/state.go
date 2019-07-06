package fsm

type State struct {
	value       string
	isAccepting bool
}

func NewState(value string) State {
	return State{
		value:       value,
		isAccepting: false,
	}
}

func NewAcceptingState(value string) State {
	return State{
		value:       value,
		isAccepting: true,
	}
}

func (s State) Accepting() bool {
	return s.isAccepting
}

func (s State) Value() string {
	return s.value
}

func (s State) Equal(other State) bool {
	return s.value == other.value
}
