package main

type Size struct {
	Rows, Cols int
}

type Index struct {
	Row, Col int
}

func (idx *Index) Neighbours() []Index {
	return []Index{
		Index{idx.Row, idx.Col + 1},
		Index{idx.Row, idx.Col - 1},
		Index{idx.Row + 1, idx.Col},
		Index{idx.Row - 1, idx.Col},
	}
}