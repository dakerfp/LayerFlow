package main

import (
	"os"
	"strings"
	"testing"
)

// TODO: read source files
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

func TestZeroSink(t *testing.T) {
	r := strings.NewReader("0!")
	tokenGrid, err := lexer(r)
	if err != nil {
		t.Fail()
	}
	program := assembleLayer(tokenGrid)
	if program == nil {
		t.Fail()
	}

	go func() {
		for i := 0; i < 10; i += 1 {
			v, ok := <-program.Output
			t.Log("zero", v, ok)
			if !ok || v != 0 {
				t.Fatal("wrong data")
			}
		}
		program.Halt <- 0
	}()

	program.Run()
}

func TestSourceSink(t *testing.T) {
	r := strings.NewReader("@!")
	tokenGrid, err := lexer(r)
	if err != nil {
		t.Fail()
	}
	program := assembleLayer(tokenGrid)
	if program == nil {
		t.Fail()
	}

	data := []Value{1, 2, 3, 4, 5}

	go func() {
		for _, v := range data {
			t.Log("send:", v)
			program.Input <- v
		}
		program.Halt <- 0
	}()

	go func() {
		var i int
		for v := range program.Output {
			t.Log("recv:", v)
			if v != data[i] {
				t.Fatal("wrong data")
			}
			i += 1
		}
		_, ok := <-program.Output
		if ok {
			t.Fatal("output should be close")
		}
	}()

	program.Run()
}
