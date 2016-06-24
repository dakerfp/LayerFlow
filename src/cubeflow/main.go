package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

var inputData = flag.String("i", "", "input data")
var debug = flag.Bool("d", false, "debug, mode")

func main() {
	flag.Parse()

	if len(flag.Args()) < 1 {
		log.Fatal(flag.Args())
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
	} else if *debug {
		log.Println(program)
	}

	go func() {
		var r io.Reader
		if *inputData != "" {
			r = strings.NewReader(*inputData)
		} else {
			r = os.Stdin
		}
		scanner := bufio.NewScanner(r)
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
			fmt.Println(os.Stdout, v)
		}
	}()

	program.Run()
	<-program.Halt
}
