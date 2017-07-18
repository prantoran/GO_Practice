package gsheets

// Row constans cell values of one row
type Row struct {
	cell []string
}

// Cell returns cell value by index
func (r *Row) Cell(i int) string {
	return r.cell[i]
}

// CellN returns length of the cells
func (r *Row) CellN() int {
	return len(r.cell)
}
