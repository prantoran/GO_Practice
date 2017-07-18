package gsheets

// Header defines the column names
type Header struct {
	pos   int
	label string
}

// Pos returns the position of the header
func (h *Header) Pos() int {
	return h.pos
}

// Label returns the label/name of the header
func (h *Header) Label() string {
	return h.label
}
