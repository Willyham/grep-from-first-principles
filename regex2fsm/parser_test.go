package regex2fsm

import (
	"regexp/syntax"
	"testing"

	"github.com/Willyham/gfp/fsm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConvert(t *testing.T) {
	parser := New()
	_, err := parser.Convert("a")
	assert.NoError(t, err)
}

func TestConvertInvalidRegex(t *testing.T) {
	parser := New()
	machine, err := parser.Convert("**invalid")
	assert.Error(t, err)
	assert.Nil(t, machine)
}

func TestParseCombined(t *testing.T) {
	parser := New()
	state := parser.stateGenerator.Next()
	tree, err := syntax.Parse("a+b*|c[d-f]", syntax.POSIX)
	require.NoError(t, err)
	transitions := parser.parseTree(state, tree)
	assert.Len(t, transitions, 7)
}

func TestParseUnsupported(t *testing.T) {
	parser := New()
	state := parser.stateGenerator.Next()
	assert.Panics(t, func() {
		tree, err := syntax.Parse("a.b", syntax.POSIX)
		require.NoError(t, err)
		parser.parseTree(state, tree)
	})
}

func TestParseLiteral(t *testing.T) {
	parser := New()
	state := parser.stateGenerator.Next()
	literal, err := syntax.Parse("a", syntax.POSIX)
	require.NoError(t, err)
	transitions := parser.parseLiteral(state, literal)
	expected := []fsm.Transition{
		{Event: "a", Source: state, NextState: fsm.NewState("1")},
	}
	assert.Equal(t, expected, transitions)
}

func TestParseLiteralMultiple(t *testing.T) {
	parser := New()
	state := parser.stateGenerator.Next()
	literal, err := syntax.Parse("abc", syntax.POSIX)
	require.NoError(t, err)
	transitions := parser.parseLiteral(state, literal)
	expected := []fsm.Transition{
		{Event: "a", Source: state, NextState: fsm.NewState("1")},
		{Event: "b", Source: fsm.NewState("1"), NextState: fsm.NewState("2")},
		{Event: "c", Source: fsm.NewState("2"), NextState: fsm.NewState("3")},
	}
	assert.Equal(t, expected, transitions)
}

func TestParseClass(t *testing.T) {
	parser := New()
	state := parser.stateGenerator.Next()
	class, err := syntax.Parse("[a-bA-B]", syntax.POSIX)
	require.NoError(t, err)
	transitions := parser.parseCharClass(state, class)
	expected := []fsm.Transition{
		{Event: "A", Source: state, NextState: fsm.NewState("1")},
		{Event: "B", Source: state, NextState: fsm.NewState("1")},
		{Event: "a", Source: state, NextState: fsm.NewState("1")},
		{Event: "b", Source: state, NextState: fsm.NewState("1")},
	}
	assert.Equal(t, expected, transitions)
}

func TestParseClassSimplifiedAlternate(t *testing.T) {
	parser := New()
	state := parser.stateGenerator.Next()
	class, err := syntax.Parse("a|b", syntax.POSIX)
	require.NoError(t, err)
	transitions := parser.parseCharClass(state, class)
	expected := []fsm.Transition{
		{Event: "a", Source: state, NextState: fsm.NewState("1")},
		{Event: "b", Source: state, NextState: fsm.NewState("1")},
	}
	assert.Equal(t, expected, transitions)
}

func TestParseAlternateSub(t *testing.T) {
	parser := New()
	state := parser.stateGenerator.Next()
	alternate, err := syntax.Parse("ab|cd", syntax.POSIX)
	require.NoError(t, err)
	transitions := parser.parseAlternate(state, alternate)
	expected := []fsm.Transition{
		{Event: "a", Source: state, NextState: fsm.NewState("1")},
		{Event: "b", Source: fsm.NewState("1"), NextState: fsm.NewState("2")},
		{Event: "c", Source: state, NextState: fsm.NewState("3")},
		{Event: "d", Source: fsm.NewState("3"), NextState: fsm.NewState("4")},
	}
	assert.Equal(t, expected, transitions)
}

func TestParsePlus(t *testing.T) {
	parser := New()
	state := parser.stateGenerator.Next()
	plus, err := syntax.Parse("a+", syntax.POSIX)
	require.NoError(t, err)
	transitions := parser.parsePlus(state, plus)
	expected := []fsm.Transition{
		{Event: "a", Source: state, NextState: fsm.NewState("1")},
		{Event: "a", Source: fsm.NewState("1"), NextState: fsm.NewState("1")},
	}
	assert.Equal(t, expected, transitions)
}

func TestParseStar(t *testing.T) {
	parser := New()
	state := parser.stateGenerator.Next()
	star, err := syntax.Parse("a*", syntax.POSIX)
	require.NoError(t, err)
	transitions := parser.parseStar(state, star)
	expected := []fsm.Transition{
		{Event: "a", Source: state, NextState: state},
	}
	assert.Equal(t, expected, transitions)
}

func TestParseConcat(t *testing.T) {
	parser := New()
	state := parser.stateGenerator.Next()
	concat, err := syntax.Parse("a*b*", syntax.POSIX)
	require.NoError(t, err)
	transitions := parser.parseConcat(state, concat)
	expected := []fsm.Transition{
		{Event: "a", Source: state, NextState: state},
		{Event: "b", Source: state, NextState: state},
	}
	assert.Equal(t, expected, transitions)
}
