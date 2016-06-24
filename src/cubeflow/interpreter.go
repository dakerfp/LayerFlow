package main

import (
	"log"
)

func (cell Cell) Clock(halt chan int) {
	for cell.Type.Exec(cell.Notify, halt) {
		log.Println("spin", cell)
	}
	close(cell.Notify)
}

func (program *Program) Run() {
	log.Println(program.Input)
	log.Println(program.Output)
	for _, cell := range program.Cells {
		log.Println(cell)
		log.Println(cell.Type)
		go cell.Clock(program.Halt)
	}
}
