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
		Size{Layers: 1, Rows: 0, Cols: 0},
		make(map[Index]rune),
	}

	scanner := bufio.NewScanner(reader)
	idx := Index{0, 0, 0}
	for scanner.Scan() {
		var r rune
		any := false
		for idx.Col, r = range scanner.Text() {
			any = true
			if r == '#' { // Comment
				break
			}
			if r == ' ' { // Whitespace
				continue
			}
			grid.Tokens[idx] = r
			if idx.Col + 1 > grid.Cols {
				grid.Cols = idx.Col + 1
			}
		}
		if !any {
			// XXX: never happens
			idx.Row = 0
			idx.Layer += 1
			grid.Layers = idx.Layer
		} else {
			if idx.Row + 1 > grid.Rows {
				grid.Rows = idx.Row + 1
			}
			idx.Row += 1
		}
	}

	return grid, nil
}
