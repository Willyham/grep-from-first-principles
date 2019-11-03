package regex2fsm

import (
	"fmt"
	"regexp/syntax"

	"github.com/Willyham/gfp/fsm"
)

type Parser struct {
	stateGenerator fsm.StateGenerator
}

func New() *Parser {
	return &Parser{
		stateGenerator: &fsm.NumericStateGenerator{},
	}
}

func (g Parser) getNextState(isAccepting bool) fsm.State {
	if isAccepting {
		return g.stateGenerator.NextAccepting()
	}
	return g.stateGenerator.Next()
}

func (g Parser) Convert(pattern string) (*fsm.StateMachine, error) {
	regexTree, err := syntax.Parse(pattern, syntax.POSIX)
	if err != nil {
		return nil, err
	}

	initialState := g.stateGenerator.Next()
	transitions := g.parseTree(initialState, regexTree, true)

	machine, err := fsm.New(
		initialState,
		transitions,
	)
	if err != nil {
		return nil, err
	}
	return machine, nil
}

func (g Parser) parseTree(currentState fsm.State, tree *syntax.Regexp, isAccepting bool) []fsm.Transition {
	switch tree.Op {
	case syntax.OpAlternate:
		return g.parseAlternate(currentState, tree, isAccepting)
	case syntax.OpLiteral:
		return g.parseLiteral(currentState, tree, isAccepting)
	case syntax.OpStar:
		return g.parseStar(currentState, tree, isAccepting)
	case syntax.OpPlus:
		return g.parsePlus(currentState, tree, isAccepting)
	case syntax.OpConcat:
		return g.parseConcat(currentState, tree, isAccepting)
	case syntax.OpCharClass:
		return g.parseCharClass(currentState, tree, isAccepting)
	default:
		panic(fmt.Sprintf("unsupported operation: %s", tree.Op))
	}
}

func (g Parser) parseAlternate(currentState fsm.State, alternate *syntax.Regexp, isAccepting bool) []fsm.Transition {
	left := g.parseTree(currentState, alternate.Sub[0], isAccepting)
	right := g.parseTree(currentState, alternate.Sub[1], isAccepting)
	return append(left, right...)
}

func (g Parser) parseLiteral(currentState fsm.State, literal *syntax.Regexp, isAccepting bool) []fsm.Transition {
	transitions := []fsm.Transition{}
	last := currentState
	for i, c := range literal.Rune {
		isLast := i == len(literal.Rune)-1
		nextState := g.getNextState(isAccepting && isLast)
		transitions = append(transitions, fsm.Transition{
			Event:     string(c),
			Source:    fsm.State(last),
			NextState: nextState,
		})
		last = nextState
	}
	return transitions
}

func (g Parser) parsePlus(currentState fsm.State, plus *syntax.Regexp, isAccepting bool) []fsm.Transition {
	midState := g.getNextState(isAccepting)
	repeatingState := g.getNextState(isAccepting)
	return []fsm.Transition{
		{Event: string(plus.Sub[0].Rune[0]), Source: currentState, NextState: midState},
		{Event: string(plus.Sub[0].Rune[0]), Source: midState, NextState: repeatingState},
		{Event: string(plus.Sub[0].Rune[0]), Source: repeatingState, NextState: repeatingState},
	}
}

func (g Parser) parseStar(currentState fsm.State, star *syntax.Regexp, isAccepting bool) []fsm.Transition {
	if isAccepting {
		currentState = currentState.MakeAccepting()
	}
	tempState := g.getNextState(isAccepting)
	return []fsm.Transition{
		{Event: string(star.Sub[0].Rune[0]), Source: currentState, NextState: tempState},
		{Event: string(star.Sub[0].Rune[0]), Source: tempState, NextState: tempState},
	}
}

func (g Parser) parseConcat(currentState fsm.State, concat *syntax.Regexp, isAccepting bool) []fsm.Transition {
	// Link current state to first state
	source := currentState
	transitions := []fsm.Transition{}
	for i, expression := range concat.Sub {
		isLast := i == len(concat.Sub)-1
		subTransitions := g.parseTree(source, expression, isAccepting && isLast)
		source = subTransitions[len(subTransitions)-1].NextState
		transitions = append(transitions, subTransitions...)
	}
	return transitions
}

func (g Parser) parseCharClass(currentState fsm.State, class *syntax.Regexp, isAccepting bool) []fsm.Transition {
	transitions := []fsm.Transition{}
	nextState := g.stateGenerator.Next()
	// Group into batches of ranges.
	// E.g [[a, b], [x, y]]
	var ranges [][]rune
	for i := 0; i < len(class.Rune); i += 2 {
		ranges = append(ranges, []rune{
			class.Rune[i],
			class.Rune[i+1],
		})
	}
	// Process each batch of two
	for _, charRange := range ranges {
		current := charRange[0]
		last := charRange[1]
		// TODO: Rather than generating every single rune, we should instead match events based on
		// functions rather than just strings.
		for current <= last {
			transitions = append(transitions, fsm.Transition{
				Event:     string(current),
				Source:    currentState,
				NextState: nextState,
			})
			current++
		}
	}
	return transitions
}
