package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
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
		if _, err := parse(tokenGrid); err != nil {
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
	program, err := parse(tokenGrid)
	if err != nil {
		t.Fail()
	}

	output := make(chan Value, 1)
	input := make(chan Value, 1)
	halt := make(chan int, 1)

	go func() {
		for i := 0; i < 10; i += 1 {
			input <- 0 // Any data
		}
		input <- 0 // Latency of 1
		close(input)
	}()

	go program.Run(input, output, halt)

	<-output // Ignore: latency of 1
	for i := 0; i < 10; i += 1 {
		v, ok := <-output
		t.Log("zero", v, ok)
		if !ok || v != 0 {
			t.Fatal("wrong data")
		}
	}
}

func TestSourceSink(t *testing.T) {
	sources := []string{"@!", "!@", "@\n!", "!\n@"}
	for _, source := range sources {
		r := strings.NewReader(source)
		tokenGrid, err := lexer(r)
		if err != nil {
			t.Fatal("lexer failed")
		}

		program, err := parse(tokenGrid)
		if err != nil {
			t.Fatal("parser failed")
		}

		output := make(chan Value, 1)
		input := make(chan Value, 1)
		halt := make(chan int, 1)
		data := []Value{1, 2, 3, 4, 5}

		go func() {
			for _, v := range data {
				t.Log("send:", v)
				input <- v
			}
			input <- 0 // Latency of 1
			close(input)
		}()

		go program.Run(input, output, halt)

		<-output // Ignore: latency of 1
		i := 0
		for v := range output {
			t.Log("recv:", v)
			if v != data[i] {
				t.Fatal("wrong data")
			}
			i += 1
		}
	}
}

func TestInterpreter(t *testing.T) {
	dirs, err := filepath.Glob("testdata/*")
	if err != nil {
		t.Fatal(err)
	}
	if len(dirs) == 0 {
		wd, _ := os.Getwd()
		t.Fatal("no tests in testdata found in: ", wd)
	}
	for _, dir := range dirs {
		fin, err := os.Open(filepath.Join(dir, "input.txt"))
		if err != nil {
			t.Fatal(err)
		}
		defer fin.Close()

		fout, err := ioutil.TempFile("", "test")
		if err != nil {
			t.Fatal(err)
		}
		defer func() {
			fout.Close()
			os.Remove(fout.Name())
		}()

		wd, _ := os.Getwd()
		cmd := exec.Command(filepath.Join(wd, "cubeflow"), "-o", fout.Name(), filepath.Join(dir, "program.cf"))
		cmd.Stdin = fin
		if err = cmd.Run(); err != nil {
			t.Fatal(err, cmd.Path)
		}

		buffer := bytes.Buffer{}
		cmd = exec.Command("diff", fout.Name(), filepath.Join(dir, "output.txt"))
		cmd.Stdout = &buffer
		if err = cmd.Run(); err != nil {
			t.Fatal(string(buffer.Bytes()))
		}
	}
}
