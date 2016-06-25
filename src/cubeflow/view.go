package main

import (
	"bufio"
	"fmt"
	"image"
	"image/color"
	"io"
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
			if _, err := rw.WriteString(r); err != nil {
				return err
			}
		}
		if _, err := rw.WriteString("\n"); err != nil {
			return err
		}
	}

	return rw.Flush()
}

func clip(v Value, clipMin, clipMax int) uint16 {
	return uint16(v << 8) // XXX: OMG
}

func DrawGrid(p *Program, clipMin, clipMax int) image.Image {
	img := image.NewGray16(image.Rect(0, 0, p.Rows, p.Cols))
	// TODO: print multiple layers
	i := 0
	for j := 0; j < p.Rows; j += 1 {
		for k := 0; k < p.Cols; k += 1 {
			cell, ok := p.Cells[Index{Layer: i, Row: j, Col: k}]
			c := color.Gray16{Y: 0}
			if ok {
				c.Y = clip(cell.Read, clipMin, clipMax)
			}
			img.Set(k, j, c)
		}
	}
	return img
}
