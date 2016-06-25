package main

func assembleLayer(grid *TokenGrid) (*Program, error) {
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
		case 'C':
			cell.Type = &Oscillator{
				Period: 1,
				Function: func(i, p uint64) Value {
					return Value(i)
				},
			}
		}
		program.Cells[idx] = cell
	}

	// Link cells
	// TODO: raise error if has no connection
	for idx, cell := range program.Cells {
		switch cell.Symbol {
		case '@':
			cell.Type.Bind(program.Input)
			continue
		case '!':
			program.Output = cell.Notify
		}

		// Try to bind all the neighbours
		for _, nidx := range idx.Neighbours(DirsPlane) {
			n, ok := program.Cells[nidx]
			if !ok {
				continue
			}
			// TODO: raise error if it is disconnected
			if err := cell.Type.Bind(n.Notify); err != nil {
				return nil, err
			}
		}
	}

	return program, nil
}
