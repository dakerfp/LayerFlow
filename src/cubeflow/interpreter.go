package main

import "log"

func (program *Program) tick() {
	// TODO: optimize evaluation order
	for _, cell := range program.Cells {
		cell.Write = cell.Type.Exec(cell.Read)
	}
	for _, cell := range program.Cells {
		if *verbose {
			log.Printf("%s -> %d %d\n", string(cell.Symbol), cell.Read, cell.Write)
		}
		cell.Read = cell.Write
	}
}

func (program *Program) Run(input, output chan Value, halt chan int) {
	// TODO: init values must be properly implemented
	for {
		var ok bool
		select {
		case program.read, ok = <-input:
			if !ok {
				close(output)
				return
			}
			program.tick()
			output <- *program.write
		case <-halt:
			close(output)
			return
		}
	}
}
