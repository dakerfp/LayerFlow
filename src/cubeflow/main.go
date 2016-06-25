package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
)

var (
	verbose = flag.Bool("v", false, "verbose mode")
	pngFilenaem = flag.String("img", "", "png filenaem")
)

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

	program, err := assembleLayer(tokenGrid)
	if err != nil {
		log.Fatal("compilation error")
	}
	if *verbose {
		err = WriteGrid(os.Stderr, program)
		if err != nil {
			log.Fatal(err)
		}
	}

	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			n, err := strconv.Atoi(scanner.Text())
			if err != nil {
				log.Fatal(err)
			}
			program.Input <- Value(n)
		}
		program.Halt <- 0
	}()

	go func() {
		for v := range program.Output {
			fmt.Fprintln(os.Stdout, v)
		}
	}()

	program.Run()
	<-program.Halt
}
