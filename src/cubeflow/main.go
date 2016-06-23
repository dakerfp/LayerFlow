package main

import (
	"fmt"
	"os"
)

func main() {
	tokenGrid, err := lexer(os.Stdin)
	if err != nil {
		os.Exit(1)
	}
	fmt.Fprintln(os.Stderr, tokenGrid)
	program := assembleLayer(tokenGrid)
	fmt.Fprintln(os.Stderr, program)
}
