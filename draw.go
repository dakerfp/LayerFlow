package main

import (
	"image"
	"image/color"
)

// Currently only the first layer is supported

type ProgramView struct {
	Program *Program
	Layer   int
}

func (view *ProgramView) ColorModel() color.Model {
	return color.RGBA64Model

}

func (view *ProgramView) Bounds() image.Rectangle {
	return image.Rectangle{
		Min: image.Point{0, 0},
		Max: image.Point{view.Program.Rows, view.Program.Cols},
	}
}

var (
	ColorRed    = color.NRGBA64{R: 255, A: 255}
	ColorGreen  = color.NRGBA64{G: 255, A: 255}
	ColorBlue   = color.NRGBA64{B: 255, A: 255}
	ColorYellow = color.NRGBA64{R: 255, G: 255, A: 255}
	ColorCyan   = color.NRGBA64{G: 255, B: 255, A: 255}
)

func (view *ProgramView) At(x, y int) color.Color {
	idx := Index{Layer: view.Layer, Row: x, Col: y}
	cell, ok := view.Program.Cells[idx]
	if !ok {
		return color.Transparent
	}

	switch cell.Type.(type) {
	case *Forward:
		switch cell.Symbol {
		case '@':
			return &ColorGreen
		case '!':
			return &ColorRed
		}
	case *Constant:
		return &ColorBlue
	case *Oscillator, *Pulse:
		return &ColorYellow
	case *BinaryOp, *UnaryOp:
		return &ColorCyan
	}

	return color.White // TODO: paint different colors for each kind of cell
}

func NewProgramView(p *Program, layer int) *ProgramView {
	return &ProgramView{p, layer}
}
