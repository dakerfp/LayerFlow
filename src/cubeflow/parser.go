package main

type Cell struct {
	Index
	Notify chan int
}

type Program struct {
	Cells map[Index]Cell
	Input, Output chan int
}

func assembleLayer(grid *TokenGrid) (*Program) {
	// TODO: support more than one input and output chan per layer
	program := &Program{Cells: make(map[Index]Cell)}
	for idx, r := range grid.Tokens {
		notify := make(chan int, 1) // TODO: optimize this
		switch r {
		case '@':
			program.Input = notify
		case '!':
			program.Output = notify
		}
		program.Cells[idx].Notify = notify
	}

	for idx, r := range grid.Tokens {
		switch r {
		case '@':
			
			
		case '!':
			
		}	
	}
	return program
}