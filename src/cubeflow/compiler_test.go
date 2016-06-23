package main

import (
	"os"
	"testing"
)

func TestCompile(t *testing.T) {
	var sourceFiles []string
	for _, source := range sourceFiles {
		f, err := os.Open(source)
		if err != nil {
			t.Fail()
		}
		defer f.Close()
		tokenGrid, err := lexer(os.Stdin)
		if err != nil {
			t.Fail()
		}
		if assembleLayer(tokenGrid) == nil {
			t.Fail()
		}
	}
}
