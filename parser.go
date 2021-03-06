package main

import (
	"fmt"
)

func parse(grid *TokenGrid) (*Program, error) {
	// TODO: support more than one input and output per layer
	program := &Program{
		Size:  grid.Size,
		Cells: make(map[Index]*Cell),
	}

	// Build cells & the layer channels
	for idx, r := range grid.Tokens {
		cell := &Cell{
			Symbol: r,
		}
		switch r {
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			n := Value(r - '0')
			cell.Type = &Constant{n}
			cell.Read = n
		case '@':
			cell.Type = &Forward{
				SourceDir: DirNone,
				SinkDir:   DirsPlane,
			}
		case '!':
			cell.Type = &Forward{
				SourceDir: DirsPlane,
				SinkDir:   DirNone,
			}
		case '^':
			cell.Type = &Forward{
				SourceDir: DirRight | DirLeft | DirDown,
				SinkDir:   DirUp,
			}
		case '<':
			cell.Type = &Forward{
				SourceDir: DirRight | DirUp | DirDown,
				SinkDir:   DirLeft,
			}
		case '>':
			cell.Type = &Forward{
				SourceDir: DirUp | DirLeft | DirDown,
				SinkDir:   DirRight,
			}
		case 'v':
			cell.Type = &Forward{
				SourceDir: DirRight | DirLeft | DirUp,
				SinkDir:   DirDown,
			}
		case 'C':
			cell.Type = &Oscillator{
				Period: 1,
				Function: func(i, p uint64) Value {
					return Value(i)
				},
			}
		case '+':
			cell.Type = &BinaryOp{
				Function: func(a, b Value) Value {
					return a + b
				},
			}
		case '*':
			cell.Type = &BinaryOp{
				Function: func(a, b Value) Value {
					return a * b
				},
			}
		case '-':
			cell.Type = &UnaryOp{
				Function: func(a Value) Value {
					return -a
				},
			}
		case 'P':
			cell.Type = &Pulse{
				Value: Value(1),
			}
		}

		program.Cells[idx] = cell
	}

	// Link cells
	// TODO: raise error if has no connection
	for idx, cell := range program.Cells {
		// Try to bind all the neighbours
		for _, dir := range Dirs(cell.Type.RequestDir()) {
			nidx, err := idx.Neighbour(dir)
			if err != nil {
				return nil, fmt.Errorf("Error in neighborhood: %v", err)
			}

			neighbour, ok := program.Cells[nidx]
			if !ok {
				continue
			}

			// Try to matching offer and request
			if InverseDir(dir)&neighbour.Type.OfferDir() == 0 {
				continue
			}

			if err := cell.Type.Bind(&neighbour.Read); err != nil {
				return nil, err
			}
		}

		switch cell.Symbol {
		case '@':
			cell.Type.Bind(&program.read)
			continue
		case '!':
			program.write = &cell.Read
		}
	}

	return program, nil
}
