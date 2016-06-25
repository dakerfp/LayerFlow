package main

import (
	"bufio"
	"io"
	"fmt"
)

func WriteGrid(w io.Writer, p *Program) error {
	rw := bufio.NewWriter(w)
	// TODO: print multiple layers
	i := 0
	rw.WriteString(fmt.Sprintf("{%v %v %v}\n",
		p.Layers, p.Rows, p.Cols))
	for j := 0; j < p.Rows; j += 1 {
		for k := 0; k < p.Cols; k += 1 {
			cell, ok := p.Cells[Index{Layer: i, Row: j, Col: k}]
			r := " "
			if ok {
				r = string(cell.Symbol)
			}
			if _, err := rw.WriteString(r); err != nil  {
				return err
			}
		}
		if _, err := rw.WriteString("\n"); err != nil {
			return err
		}
	}

	return rw.Flush()
}
