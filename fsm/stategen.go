package fsm

import (
	"strconv"
	"sync"
)

type StateGenerator interface {
	Next() State
	NextAccepting() State
}

type NumericStateGenerator struct {
	counter int
	sync.Mutex
}

func (g *NumericStateGenerator) Next() State {
	g.Lock()
	defer g.Unlock()
	state := NewState(strconv.Itoa(g.counter))
	g.counter++
	return state
}

func (g *NumericStateGenerator) NextAccepting() State {
	g.Lock()
	defer g.Unlock()
	state := NewAcceptingState(strconv.Itoa(g.counter))
	g.counter++
	return state
}
