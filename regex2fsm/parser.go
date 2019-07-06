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

func (g Parser) Convert(pattern string) (*fsm.StateMachine, error) {
	regexTree, err := syntax.Parse(pattern, syntax.POSIX)
	if err != nil {
		return nil, err
	}

	initialState := g.stateGenerator.Next()
	transitions := g.parseTree(initialState, regexTree)
	machine, err := fsm.New(
		initialState,
		transitions,
	)
	if err != nil {
		return nil, err
	}
	return machine, nil
}

func (g Parser) parseTree(currentState fsm.State, tree *syntax.Regexp) []fsm.Transition {
	switch tree.Op {
	case syntax.OpAlternate:
		return g.parseAlternate(currentState, tree)
	case syntax.OpLiteral:
		return g.parseLiteral(currentState, tree)
	case syntax.OpStar:
		return g.parseStar(currentState, tree)
	case syntax.OpPlus:
		return g.parsePlus(currentState, tree)
	case syntax.OpConcat:
		return g.parseConcat(currentState, tree)
	case syntax.OpCharClass:
		return g.parseCharClass(currentState, tree)
	default:
		panic(fmt.Sprintf("unsuported operation: %s", tree.Op))
	}
}

func (g Parser) parseAlternate(currentState fsm.State, alternate *syntax.Regexp) []fsm.Transition {
	left := g.parseTree(currentState, alternate.Sub[0])
	right := g.parseTree(currentState, alternate.Sub[1])
	return append(left, right...)
}

func (g Parser) parseLiteral(currentState fsm.State, literal *syntax.Regexp) []fsm.Transition {
	transitions := []fsm.Transition{}
	last := currentState
	for _, c := range literal.Rune {
		nextState := g.stateGenerator.Next()
		transitions = append(transitions, fsm.Transition{
			Event:     string(c),
			Source:    fsm.State(last),
			NextState: nextState,
		})
		last = nextState
	}
	return transitions
}

func (g Parser) parsePlus(currentState fsm.State, plus *syntax.Regexp) []fsm.Transition {
	tempState := g.stateGenerator.Next()
	return []fsm.Transition{
		{Event: string(plus.Sub[0].Rune[0]), Source: currentState, NextState: tempState},
		{Event: string(plus.Sub[0].Rune[0]), Source: tempState, NextState: tempState},
	}
}

func (g Parser) parseStar(currentState fsm.State, star *syntax.Regexp) []fsm.Transition {
	return []fsm.Transition{
		{Event: string(star.Sub[0].Rune[0]), Source: currentState, NextState: currentState},
	}
}

func (g Parser) parseConcat(currentState fsm.State, concat *syntax.Regexp) []fsm.Transition {
	// Link current state to first state
	source := currentState
	transitions := []fsm.Transition{}
	for _, expression := range concat.Sub {
		subTransitions := g.parseTree(source, expression)
		source = subTransitions[len(subTransitions)-1].NextState
		transitions = append(transitions, subTransitions...)
	}
	return transitions
}

func (g Parser) parseCharClass(currentState fsm.State, class *syntax.Regexp) []fsm.Transition {
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
