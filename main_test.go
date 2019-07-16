package main

import (
	"bufio"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/Willyham/gfp/regex2fsm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func BenchmarkSearchString(b *testing.B) {
	converter := regex2fsm.New()
	machine, err := converter.Convert("go+")
	require.NoError(b, err)
	testString := "gracious great balls of fire goodness"
	testStringTokens := strings.Split(testString, "")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := machine.Run(testStringTokens)
		assert.True(b, result)
		machine.Reset()
	}
}

func BenchmarkDictionarySearch(b *testing.B) {
	file, err := os.Open("./words.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	tokens := [][]string{}
	for scanner.Scan() {
		line := scanner.Text()
		tokens = append(tokens, strings.Split(line, ""))
	}

	converter := regex2fsm.New()
	machine, err := converter.Convert("go")
	require.NoError(b, err)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, t := range tokens {
			machine.Run(t)
			machine.Reset()
		}
	}
}
