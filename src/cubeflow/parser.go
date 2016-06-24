package main

type Value int64

type Type interface {
	Exec(notify chan Value, halt chan int) bool
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
	Halt          chan int
}

type Forward struct {
	Input chan Value
}

func (s *Forward) Exec(notify chan Value, halt chan int) bool {
	select {
	case v, ok := <-s.Input:
		if !ok {
			return false
		}
		select {
		case <-halt:
			return false
		case notify <- v:
			return true
		}
	case <-halt:
		return false
	}
}

type Constant struct {
	Value
}

func (c *Constant) Exec(notify chan Value, halt chan int) bool {
	select {
	case <-halt:
		return false
	case notify <- c.Value:
		return true
	}
}

func assembleLayer(grid *TokenGrid) *Program {
	// TODO: support more than one input and output chan per layer
	program := &Program{
		Size:   grid.Size,
		Cells:  make(map[Index]Cell),
		Input:  make(chan Value, 1),
		Halt:   make(chan int, 1),
		Output: nil,
	}

	// Build cells & the layer channels
	for idx, r := range grid.Tokens {
		cell := Cell{
			Index:  idx,
			Notify: make(chan Value, 1), // TODO: allow buffering
			Symbol: r,
		}
		switch r {
		case '0':
			cell.Type = &Constant{0}
		case '@':
			cell.Type = &Forward{}
		case '!':
			cell.Type = &Forward{}
		}
		program.Cells[idx] = cell
	}

	// Link cells
	// TODO: raise error if has no connection
	for idx, cell := range program.Cells {
		switch cell.Symbol {
		case '@':
			cell.Notify = program.Input
		case '!':
			program.Output = cell.Notify
			sink := cell.Type.(*Forward)
			for _, nidx := range idx.Neighbours() {
				if n, ok := program.Cells[nidx]; ok {
					sink.Input = n.Notify
					break // TODO: check for errors
				}
			}
		}
	}

	return program
}
