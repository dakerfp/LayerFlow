package main

type Value int
type Type interface {
}

type Cell struct {
	Index
	Symbol rune
	Notify chan Value
	Type   Type
}

type Program struct {
	Size
	Cells         map[Index]Cell
	Input, Output chan Value
}

type Source struct {
}

type Sink struct {
	Inputs []chan Value
}

func assembleLayer(grid *TokenGrid) *Program {
	// TODO: support more than one input and output chan per layer
	program := &Program{
		Size:  grid.Size,
		Cells: make(map[Index]Cell),
	}

	// Build cells & the layer channels
	for idx, r := range grid.Tokens {
		cell := Cell{
			Index:  idx,
			Notify: make(chan Value, 1), // TODO: allow buffering
			Symbol: r,
		}
		switch r {
		case '@':
			cell.Type = &Source{}
		case '!':
			cell.Type = &Sink{}
		}
		program.Cells[idx] = cell
	}

	// Link cells
	// TODO: raise error if has no connection
	for idx, cell := range program.Cells {
		switch cell.Symbol {
		case '@':
			program.Input = cell.Notify
		case '!':
			program.Output = cell.Notify
			sink := cell.Type.(*Sink)
			for _, nidx := range idx.Neighbours() {
				if n, ok := program.Cells[nidx]; ok {
					sink.Inputs = append(sink.Inputs, n.Notify)
				}
			}
		}
	}

	return program
}
