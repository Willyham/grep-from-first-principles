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
	transitions := parser.parseTree(state, tree, false)
	assert.Len(t, transitions, 9)
}

func TestParseUnsupported(t *testing.T) {
	parser := New()
	state := parser.stateGenerator.Next()
	assert.Panics(t, func() {
		tree, err := syntax.Parse("a.b", syntax.POSIX)
		require.NoError(t, err)
		parser.parseTree(state, tree, false)
	})
}

func TestParseLiteral(t *testing.T) {
	parser := New()
	state := parser.stateGenerator.Next()
	literal, err := syntax.Parse("a", syntax.POSIX)
	require.NoError(t, err)
	transitions := parser.parseLiteral(state, literal, false)
	expected := []fsm.Transition{
		{Event: "a", Source: state, NextState: fsm.NewState("1")},
	}
	assert.Equal(t, expected, transitions)
}

func TestParseLiteralAccepting(t *testing.T) {
	parser := New()
	state := parser.stateGenerator.Next()
	literal, err := syntax.Parse("a", syntax.POSIX)
	require.NoError(t, err)
	transitions := parser.parseLiteral(state, literal, true)
	expected := []fsm.Transition{
		{Event: "a", Source: state, NextState: fsm.NewAcceptingState("1")},
	}
	assert.Equal(t, expected, transitions)
}

func TestParseLiteralMultiple(t *testing.T) {
	parser := New()
	state := parser.stateGenerator.Next()
	literal, err := syntax.Parse("abc", syntax.POSIX)
	require.NoError(t, err)
	transitions := parser.parseLiteral(state, literal, false)
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
	transitions := parser.parseCharClass(state, class, false)
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
	transitions := parser.parseCharClass(state, class, false)
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
	transitions := parser.parseAlternate(state, alternate, false)
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
	initialState := parser.stateGenerator.Next()
	plus, err := syntax.Parse("a+", syntax.POSIX)
	require.NoError(t, err)
	transitions := parser.parsePlus(initialState, plus, false)
	expected := []fsm.Transition{
		{Event: "a", Source: initialState, NextState: fsm.NewState("1")},
		{Event: "a", Source: fsm.NewState("1"), NextState: fsm.NewState("2")},
		{Event: "a", Source: fsm.NewState("2"), NextState: fsm.NewState("2")},
	}
	assert.Equal(t, expected, transitions)
}

func TestParseStar(t *testing.T) {
	parser := New()
	initialState := parser.stateGenerator.Next()
	star, err := syntax.Parse("a*", syntax.POSIX)
	require.NoError(t, err)
	transitions := parser.parseStar(initialState, star, false)
	expected := []fsm.Transition{
		{Event: "a", Source: initialState, NextState: fsm.NewState("1")},
		{Event: "a", Source: fsm.NewState("1"), NextState: fsm.NewState("1")},
	}
	assert.Equal(t, expected, transitions)
}

func TestParseConcat(t *testing.T) {
	parser := New()
	state := parser.stateGenerator.Next()
	concat, err := syntax.Parse("a*b*", syntax.POSIX)
	require.NoError(t, err)
	transitions := parser.parseConcat(state, concat, false)
	expected := []fsm.Transition{
		{Event: "a", Source: state, NextState: fsm.NewState("1")},
		{Event: "a", Source: fsm.NewState("1"), NextState: fsm.NewState("1")},
		{Event: "b", Source: fsm.NewState("1"), NextState: fsm.NewState("2")},
		{Event: "b", Source: fsm.NewState("2"), NextState: fsm.NewState("2")},
	}
	assert.Equal(t, expected, transitions)
}

func BenchmarkParse(b *testing.B) {
	for i := 0; i < b.N; i++ {
		parser := New()
		parser.Convert("a*b+|cd[e-g]+")
	}
}
