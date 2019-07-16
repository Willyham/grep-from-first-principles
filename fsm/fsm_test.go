package fsm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var openState = NewAcceptingState("open")
var closedState = NewAcceptingState("closed")
var smashedState = NewState("smashed")

// Events
var open = "open"
var close = "close"
var smash = "smash"
var invalid = "foo"

func TestSimpleMachine(t *testing.T) {
	initial := openState
	cases := []struct {
		name     string
		input    []string
		expected bool
	}{
		// {"Simple case", []string{close, open, close}, true},
		{"Not in alphabet", []string{close, invalid, open}, false},
		// {"No valid transition", []string{close, close}, false},
		// {"No valid transition (after some valid)", []string{close, open, open}, false},
		// {"Not accepted final state", []string{close, smash}, false},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			machine, err := New(initial, []Transition{
				{Event: "open", Source: closedState, NextState: openState},
				{Event: "close", Source: openState, NextState: closedState},
				{Event: "smash", Source: closedState, NextState: smashedState},
			})
			assert.NoError(t, err)
			result := machine.Run(testCase.input)
			assert.Equal(t, testCase.expected, result)
		})
	}
}

func TestToGraphViz(t *testing.T) {
	initial := openState
	machine, err := New(initial, []Transition{
		{Event: "open", Source: closedState, NextState: openState},
		{Event: "close", Source: openState, NextState: closedState},
		{Event: "smash", Source: closedState, NextState: smashedState},
	})
	assert.NoError(t, err)
	output := machine.ToGraphViz()

	expected := `digraph FSM {
	closed->open[ label=open ];
	open->closed[ label=close ];
	closed->smashed[ label=smash ];
	closed [ shape=doublecircle ];
	open [ shape=doublecircle ];
	smashed;

}
`
	assert.Equal(t, expected, output)
}
