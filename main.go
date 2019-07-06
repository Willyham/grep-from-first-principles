package main

import (
	"fmt"

	"github.com/Willyham/gfp/regex2fsm"
)

func main() {
	converter := regex2fsm.New()
	machine, err := converter.Convert("a*b|cd+")
	if err != nil {
		panic(err)
	}

	fmt.Printf(machine.ToGraphViz())
	result := machine.Run([]string{"a", "a", "a"})
	fmt.Printf("Result: %t\n", result)
}
