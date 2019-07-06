package main

import (
	"fmt"
	"regexp/syntax"

	"github.com/Willyham/gfp/fsm"
)

func main() {
	generator := NewFSMGenerator()
	machine, err := generator.RegexToFSM("a*")
	if err != nil {
		panic(err)
	}

	fmt.Printf(machine.ToGraphViz())
	result := machine.Run([]string{"a", "a", "b"})
	fmt.Printf("Result: %t\n", result)
}

type FSMGenerator struct {
	stateGenerator fsm.StateGenerator
}

func NewFSMGenerator() *FSMGenerator {
	return &FSMGenerator{
		stateGenerator: &fsm.NumericStateGenerator{},
	}
}

func (g FSMGenerator) RegexToFSM(pattern string) (*fsm.StateMachine, error) {
	regexTree, err := syntax.Parse(pattern, syntax.POSIX)
	if err != nil {
		return nil, err
	}

	initialState := g.stateGenerator.Next()

	walkTree(regexTree, "")
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

func walkTree(tree *syntax.Regexp, pad string) ([]fsm.Transition, error) {
	fmt.Println(pad, tree.Op, tree.Rune)
	for _, node := range tree.Sub {
		_, err := walkTree(node, pad+"-")
		if err != nil {
			return nil, err
		}
	}
	return []fsm.Transition{}, nil
}

func (g FSMGenerator) parseTree(currentState fsm.State, tree *syntax.Regexp) []fsm.Transition {
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
	default:
		panic(fmt.Sprintf("unsuported operation: %s", tree.Op))
	}
}

func (g FSMGenerator) parseAlternate(currentState fsm.State, alternate *syntax.Regexp) []fsm.Transition {
	left := g.parseTree(currentState, alternate.Sub[0])
	right := g.parseTree(currentState, alternate.Sub[1])
	return append(left, right...)
}

func (g FSMGenerator) parseLiteral(currentState fsm.State, literal *syntax.Regexp) []fsm.Transition {
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

func (g FSMGenerator) parseStar(currentState fsm.State, star *syntax.Regexp) []fsm.Transition {
	tempState := g.stateGenerator.Next()
	return []fsm.Transition{
		{Event: string(star.Sub[0].Rune[0]), Source: currentState, NextState: tempState},
		{Event: string(star.Sub[0].Rune[0]), Source: tempState, NextState: tempState},
		// {Event: "", Source: currentState, NextState: tempState},
	}
}

func (g FSMGenerator) parsePlus(currentState fsm.State, plus *syntax.Regexp) []fsm.Transition {
	tempState := g.stateGenerator.Next()
	star := g.parseStar(tempState, plus)
	transitions := append(
		[]fsm.Transition{
			{Event: string(plus.Sub[0].Rune[0]), Source: currentState, NextState: tempState},
		},
		star...,
	)
	return transitions
}

func (g FSMGenerator) parseConcat(currentState fsm.State, concat *syntax.Regexp) []fsm.Transition {
	// Link current state to first state
	source := currentState
	transitions := []fsm.Transition{}
	for _, expression := range concat.Sub {
		subTransitions := g.parseTree(source, expression)
		source = subTransitions[len(subTransitions)-1].Source
		transitions = append(transitions, subTransitions...)
	}
	return transitions
}
