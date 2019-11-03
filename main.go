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
	pattern := os.Args[1]
	filename := os.Args[2]

	converter := regex2fsm.New()
	machine, err := converter.Convert(pattern)
	if err != nil {
		log.Fatal(err)
	}
	searchInFile(filename, machine)
}

func searchInFile(filename string, machine *fsm.StateMachine) {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		tokens := strings.Split(line, "")
		result := machine.Run(tokens)
		if result {
			fmt.Println(line)
		}
		machine.Reset()
	}

	err = scanner.Err()
	if err != nil {
		log.Fatal(err)
	}
}
