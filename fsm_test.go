package fsm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSimpleMachine(t *testing.T) {
	openState := State("open")
	closedState := State("closed")
	smashedState := State("smashed")

	// Events
	open := "open"
	close := "close"
	smash := "smash"
	invalid := "foo"

	initial := openState
	finalStates := []State{
		openState,
		closedState,
	}
	cases := []struct {
		name     string
		input    []string
		expected bool
	}{
		{"Simple case", []string{close, open, close}, true},
		{"Not in alphabet", []string{close, invalid, open}, false},
		{"No valid transition", []string{close, close}, false},
		{"No valid transition (after some valid)", []string{close, open, open}, false},
		{"Not accepted final state", []string{close, smash}, false},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			machine, err := New(initial, finalStates, []transition{
				{name: "open", source: closedState, nextState: openState},
				{name: "close", source: openState, nextState: closedState},
				{name: "smash", source: closedState, nextState: smashedState},
			})
			assert.NoError(t, err)
			result := machine.Run(testCase.input)
			assert.Equal(t, testCase.expected, result)
		})
	}

}
