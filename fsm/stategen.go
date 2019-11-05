package fsm

import (
	"strconv"
	"sync"
)

// StateGenerator allows generation of new states.
type StateGenerator interface {
	Next() State
	NextAccepting() State
}

// NumericStateGenerator creates states with incrementing integer ids.
type NumericStateGenerator struct {
	counter int
	sync.Mutex
}

// Next generates a newm non-accepting state.
func (g *NumericStateGenerator) Next() State {
	g.Lock()
	defer g.Unlock()
	state := NewState(strconv.Itoa(g.counter))
	g.counter++
	return state
}

// NextAccepting generates a new accepting state.
func (g *NumericStateGenerator) NextAccepting() State {
	g.Lock()
	defer g.Unlock()
	state := NewAcceptingState(strconv.Itoa(g.counter))
	g.counter++
	return state
}
