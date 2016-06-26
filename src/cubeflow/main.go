package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
)

var (
	verbose     = flag.Bool("v", false, "verbose mode")
	pngFilenaem = flag.String("img", "", "png filenaem")
	latency     = flag.Int("lat", 1, "latency")
	head        = flag.Int("n", -1, "prints n first results after latency")
)

func ReadInts(r io.Reader, output chan Value) error {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		n, err := strconv.Atoi(scanner.Text())
		if err != nil {
			return err
		}
		output <- Value(n)
	}
	return nil
}

func main() {
	flag.Parse()

	if len(flag.Args()) < 1 {
		log.Fatal("a script is required")
	}

	f, err := os.Open(flag.Arg(0))
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	tokenGrid, err := lexer(f)
	if err != nil {
		log.Fatal(err)
	}

	program, err := parse(tokenGrid)
	if err != nil {
		log.Fatal("compilation error", err)
	}
	if *verbose {
		err = WriteGrid(os.Stderr, program)
		if err != nil {
			log.Fatal(err)
		}
	}

	input := make(chan Value, 1)
	output := make(chan Value, 1)
	halt := make(chan int, 1)

	go func() {
		if err := ReadInts(os.Stdin, input); err != nil {
			log.Fatal(err)
		}
		for i := 0; i < *latency; i += 1 {
			input <- Value(0)
		}
		close(input)

	}()

	go program.Run(input, output, halt)

	for i := 0; i < *latency; i += 1 {
		<-output
	}

	for i := 0; *head < 0 || i < *head; i += 1 {
		v, ok := <-output
		if !ok {
			break
		}
		fmt.Fprintln(os.Stdout, v)
	}
}
