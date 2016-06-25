package main

func parse(grid *TokenGrid) (*Program, error) {
	// TODO: support more than one input and output per layer
	program := &Program{
		Size:  grid.Size,
		Cells: make(map[Index]*Cell),
	}

	// Build cells & the layer channels
	for idx, r := range grid.Tokens {
		cell := &Cell{
			Index:  idx,
			Symbol: r,
		}
		switch r {
		case '0':
			cell.Type = &Constant{0}
		case '1':
			cell.Type = &Constant{1}
			cell.Read = 1
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
			cell.Type.Bind(&program.read)
			continue
		case '!':
			program.write = &cell.Read
		}

		// Try to bind all the neighbours
		for _, nidx := range idx.Neighbours(DirsPlane) {
			n, ok := program.Cells[nidx]
			if !ok {
				continue
			}
			// TODO: raise error if it is disconnected
			if err := cell.Type.Bind(&n.Read); err != nil {
				return nil, err
			}
		}
	}

	return program, nil
}
