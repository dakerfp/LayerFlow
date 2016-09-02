package main

import (
	"errors"
)

var ErrTooManyDirs = errors.New("Too many directions")

const (
	DirUp = 1 << iota
	DirDown
	DirLeft
	DirRight
	DirTop
	DirBottom
	DirsPlane = DirUp | DirDown | DirLeft | DirRight
	DirsAll   = DirsPlane | DirTop | DirBottom
	DirNone   = 0
)

type Dir int

type Size struct {
	Layers, Rows, Cols int
}

type Index struct {
	Layer, Row, Col int
}

func InverseDir(dir Dir) Dir {
	switch dir {
	case DirUp:
		return DirDown
	case DirDown:
		return DirUp
	case DirLeft:
		return DirRight
	case DirRight:
		return DirLeft
	case DirTop:
		return DirBottom
	case DirBottom:
		return DirTop
	}
	return DirNone
}

func (idx *Index) Neighbour(dir Dir) (Index, error) {
	switch dir {
	case DirUp:
		return Index{idx.Layer, idx.Row - 1, idx.Col}, nil
	case DirDown:
		return Index{idx.Layer, idx.Row + 1, idx.Col}, nil
	case DirLeft:
		return Index{idx.Layer, idx.Row, idx.Col - 1}, nil
	case DirRight:
		return Index{idx.Layer, idx.Row, idx.Col + 1}, nil
	case DirTop:
		return Index{idx.Layer - 1, idx.Row, idx.Col}, nil
	case DirBottom:
		return Index{idx.Layer + 1, idx.Row, idx.Col}, nil
	}
	return Index{}, ErrTooManyDirs
}

func Dirs(dir Dir) []Dir {
	dirs := []Dir{}
	for _, d := range []Dir{DirUp, DirDown, DirLeft, DirRight, DirTop, DirBottom} {
		if dir&d != 0 {
			dirs = append(dirs, d)
		}
	}
	return dirs
}
