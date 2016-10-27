package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"image/png"
)

var (
	verbose        = flag.Bool("v", false, "verbose mode")
	outputFilename = flag.String("o", "", "output into file")
	latency        = flag.Int("lat", 1, "latency")
	head           = flag.Int("n", -1, "prints n first results after latency")
	drawFilename   = flag.String("draw", "", "draw program into a png file")
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

	// Parse
	tokenGrid, err := lexer(f)
	if err != nil {
		log.Fatal(err)
	}

	// Compile program
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

	// Draw routine
	if *drawFilename != "" {
		img := NewProgramView(program, 0) // Layer 0
		out, err := os.OpenFile(*drawFilename, os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			log.Fatal(err)	
		}
		defer out.Close()

		err = png.Encode(out, img)
		if err != nil {
			log.Fatal(err)
		}
	}

	input := make(chan Value, 1)
	output := make(chan Value, 1)
	halt := make(chan int, 1)

	// Reading values from input
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

	out := os.Stdout
	if *outputFilename != "" {
		out, err = os.OpenFile(*outputFilename, os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			log.Fatal(err)
		}
		defer out.Close()
	}

	for i := 0; *head < 0 || i < *head; i += 1 {
		v, ok := <-output
		if !ok {
			break
		}
		fmt.Fprintln(out, v)
	}
}
