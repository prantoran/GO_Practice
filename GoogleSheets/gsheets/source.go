package gsheets

import (
	"errors"

	"git.meghdut.io/meghdut/webapi/importer/parser"
)

// Source contains the row and the reader
type Source struct {
	row        Row
	currentRow int
	header     bool
	headers    []Header
	err        error
	sheet      *GSheet
}

// Next function returns bool to indicate whether another record exists and updates the row field
func (s *Source) Next() bool {

	if s.currentRow >= len(s.sheet.Rows) {
		return false
	}

	row := s.sheet.Rows[s.currentRow]
	s.currentRow++

	var cells []string

	for _, cell := range row.Cells {
		c, err := cell.String()
		if err != nil {
			return false
		}
		cells = append(cells, c)
	}

	s.row.cell = cells
	return true
}

// Err function returns the error
func (s *Source) Err() error {
	return s.err
}

// Row function returns the row
func (s *Source) Row() parser.Row {
	return &s.row
}

// Headers returns all the headers
func (s *Source) Headers() []parser.Header {
	headers := []parser.Header{}

	for _, header := range s.headers {
		headers = append(headers, &header)
	}

	return headers
}

// NewSource returns a new source object parser.Source
func NewSource(spreadsheetID string, readRange string, opts ...SourceOption) parser.Source {
	source := &Source{}

	gsheet := NewGSheet(spreadsheetID, readRange)
	err := gsheet.Connect()
	err = gsheet.Read()

	//file, err := xlsx.OpenReaderAt(reader, fileSize)

	if err != nil {
		source.err = err
		return source
	}

	if len(file.Sheets) == 0 {
		source.err = errors.New("No sheet found!")
		return source
	}

	// get only first sheet
	//source.sheet = file.Sheets[0]
	source.sheet = gsheet

	//set current row to 0
	source.currentRow = 0

	for _, opt := range opts {
		opt.Apply(source)
	}

	if source.header {
		if source.Next() {
			row := source.Row()
			for i := 0; i < row.CellN(); i++ {
				header := Header{}
				header.pos = i
				header.label = row.Cell(i)
				source.headers = append(source.headers, header)
			}
		}
	}

	return source
}
