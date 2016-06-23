package main

import (
	"os"
	"fmt"
)

func main() {
	tokenGrid, err := lexer(os.Stdin)
	if err != nil {
		os.Exit(1)
	}
	fmt.Fprintln(os.Stderr, tokenGrid)
}
