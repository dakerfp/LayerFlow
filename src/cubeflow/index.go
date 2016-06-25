package main

const (
	DirUp = 1 << iota
	DirDown
	DirLeft
	DirRight
	DirTop
	DirBottom
	DirsPlane = DirUp | DirDown | DirLeft | DirRight
	DirsAll   = DirsPlane | DirTop | DirBottom
)

type Dir int

type Size struct {
	Layers, Rows, Cols int
}

type Index struct {
	Layer, Row, Col int
}

func (idx *Index) Neighbour(dir Dir) (Index, bool) {
	switch dir {
	case DirUp:
		return Index{idx.Layer, idx.Row - 1, idx.Col}, true
	case DirDown:
		return Index{idx.Layer, idx.Row + 1, idx.Col}, true
	case DirLeft:
		return Index{idx.Layer, idx.Row, idx.Col - 1}, true
	case DirRight:
		return Index{idx.Layer, idx.Row, idx.Col + 1}, true
	case DirTop:
		return Index{idx.Layer - 1, idx.Row, idx.Col}, true
	case DirBottom:
		return Index{idx.Layer + 1, idx.Row, idx.Col}, true
	}
	return *idx, false
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

func (idx *Index) Neighbours(dir Dir) []Index {
	dirs := Dirs(dir)
	idxs := make([]Index, len(dirs))
	for i, d := range dirs {
		idxs[i], _ = idx.Neighbour(d)
	}
	return idxs
}
