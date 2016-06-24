package main

type Size struct {
	Rows, Cols, Layers int
}

type Index struct {
	Row, Col, Layer int
}

func (idx *Index) Neighbours() []Index {
	return []Index{
		Index{idx.Row, idx.Col + 1, idx.Layer},
		Index{idx.Row, idx.Col - 1, idx.Layer},
		Index{idx.Row + 1, idx.Col, idx.Layer},
		Index{idx.Row - 1, idx.Col, idx.Layer},
	}
}
