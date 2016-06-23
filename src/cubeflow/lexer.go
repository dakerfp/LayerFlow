package main

import (
	"bufio"
	"io"
)

type TokenGrid struct {
	Size
	Tokens map[Index]rune
}

func lexer(reader io.Reader) (*TokenGrid, error) {
	// TODO: currently lexer only reads a single 2D grid
	grid := &TokenGrid{
		Size{0, 0},
		make(map[Index]rune),
	}

	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		for i, r := range scanner.Text() {
			if r == '#' { // Comment
				break
			}
			if r == ' ' { // Whitespace
				continue
			}
			grid.Tokens[Index{Row: grid.Rows, Col: i}] = r
			if i > grid.Cols {
				grid.Cols = i
			}
		}
		grid.Rows += 1
	}
	return grid, nil
}
