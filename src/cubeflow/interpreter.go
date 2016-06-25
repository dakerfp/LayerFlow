package main

import (
	"log"
)

func (cell Cell) Clock(halt chan int) {
	for cell.Type.Exec(cell.Notify, halt) {
		if *verbose {
			log.Println("spin", cell)
		}
	}
	close(cell.Notify)
}

func (program *Program) Run() {
	if *verbose {
		log.Println(program.Input)
		log.Println(program.Output)
	}
	for _, cell := range program.Cells {
		if *verbose {
			log.Println(cell)
			log.Println(cell.Type)
		}
		go cell.Clock(program.Halt)
	}
}
