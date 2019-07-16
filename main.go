package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/Willyham/gfp/fsm"

	"github.com/Willyham/gfp/regex2fsm"
)

func main() {
	converter := regex2fsm.New()
	machine, err := converter.Convert("[a-b]+")
	if err != nil {
		panic(err)
	}

	fmt.Printf(machine.ToGraphViz())
	result := machine.Run([]string{"g", "o", "o", "d"})
	fmt.Printf("Result: %t\n", result)
	// searchInFile(machine)
}

func searchInFile(machine *fsm.StateMachine) {
	file, err := os.Open("./words.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		tokens := strings.Split(line, "")
		// fmt.Println(tokens)
		result := machine.Run(tokens)
		if result {
			fmt.Printf("Found match in line: %s\n", line)
		}
		machine.Reset()
	}

	err = scanner.Err()
	if err != nil {
		log.Fatal(err)
	}
}
