package csv

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
)

// Reader is a csv parser. It works well with broken csv files, that don't meet RFC requirements.
// It supports unquoted quotes along with multi row records.
type Reader struct {
	// A csv field separator, ',' is default separator
	Separator rune
	// If this flag is enabled, parser will check that quotes must be quoted.
	// Disabled as default.
	QuotedQuotes bool
	// You can set the desired field count. If it is not set, parser will take this value from
	// the first row.
	FieldsCount int
	// If this flag is true, parser will return an error if some lines field count differs from
	// FieldsCount value. Default true.
	FieldsFixed bool
	// OneRowRecord if true, reader will consider that one record must be on a single row
	OneRowRecord bool

	input        *bufio.Reader
	lineBuffer   bytes.Buffer
	fieldIndexes []int

	rowCount        int
	colCount        int
	currentRowIndex int

	parsingErr error
	values     []string
}

// NewReader is a constructor for a Reader.
func NewReader(r io.Reader) *Reader {
	return &Reader{
		Separator:   ',',
		FieldsFixed: true,

		input:    bufio.NewReader(r),
		rowCount: 1,
	}
}

// Err returns an error, if there was one.
func (r *Reader) Err() error {
	if r.parsingErr == io.EOF {
		return nil
	}

	return r.parsingErr
}

// Values return last successfully parsed record values.
func (r *Reader) Values() []string {
	return r.values
}

// Next reads the next record and returns if the operation was successful.
// TODO https://en.wikipedia.org/wiki/Byte_order_mark handle BOM in the beginig of the file
func (r *Reader) Next() bool {
	if r.parsingErr != nil {
		return false
	}

	r.values, r.parsingErr = r.getRecord()

	return len(r.values) > 0

}

// CurrentRowIndex returns current line number
func (r *Reader) CurrentRowIndex() int {
	return r.currentRowIndex
}

// GetNext reads the next row and returns its values.
func (r *Reader) GetNext() []string {
	r.Next()

	return r.Values()
}

func (r *Reader) getRecord() ([]string, error) {
	err := r.parseRecord()
	if err != nil && err != io.EOF || len(r.fieldIndexes) == 0 {
		return nil, err
	}

	line := r.lineBuffer.String()
	fieldCount := len(r.fieldIndexes)
	fields := make([]string, fieldCount)

	for i, idx := range r.fieldIndexes {
		if i == fieldCount-1 {
			fields[i] = line[idx:]
		} else {
			fields[i] = line[idx:r.fieldIndexes[i+1]]
		}
	}

	if r.FieldsCount == 0 {
		r.FieldsCount = fieldCount
	} else if r.FieldsFixed && r.FieldsCount != fieldCount {
		return nil, fmt.Errorf("wrong column count, row index %d", r.currentRowIndex)
	}

	return fields, err
}

func (r *Reader) parseRecord() error {
	r.lineBuffer.Reset()
	r.fieldIndexes = r.fieldIndexes[:0]
	r.currentRowIndex = r.rowCount

	for {
		r1, err := r.readRune()
		if err != nil {
			if len(r.fieldIndexes) > 0 {
				r.setNextFieldIndex()
			}

			return err
		}

		switch r1 {
		case '\n':
			r.setNextFieldIndex()

			return nil

		case r.Separator:
			r.setNextFieldIndex()

		default:
			stop, err := r.getField(r1)
			if err != nil {
				return err
			}
			if stop {
				return nil
			}
		}
	}
}

func (r *Reader) setNextFieldIndex() {
	r.fieldIndexes = append(r.fieldIndexes, r.lineBuffer.Len())
}

func (r *Reader) getField(r1 rune) (bool, error) {
	r.setNextFieldIndex()

	if r.OneRowRecord {
		return r.getUnquotedField(r1)
	}

	switch r1 {
	case '"', '\'':
		return r.getQuotedField(r1)
	default:
		return r.getUnquotedField(r1)
	}
}

func (r *Reader) getUnquotedField(r1 rune) (bool, error) {
	var err error
	r.lineBuffer.WriteRune(r1)
	for {
		r1, err = r.readRune()
		if err != nil {
			return true, err
		}

		if r1 == '\n' || r1 == r.Separator {
			return r1 == '\n', nil
		}

		r.lineBuffer.WriteRune(r1)
	}
}

func (r *Reader) getQuotedField(pair rune) (bool, error) {
	var r1, r2 rune
	var err error

	for {
		r1 = r2
		r2, err = r.readRune()
		if err != nil {
			return true, err
		}

		if r2 == '\n' || r2 == r.Separator {
			if r1 == pair {
				return r2 == '\n', nil
			}
		}
		if r.QuotedQuotes {
			if r1 == pair && r2 != pair {
				return true, fmt.Errorf("unquoted quote on %d:%d", r.currentRowIndex, r.colCount-1)
			}

			if r1 == pair && r2 == pair {
				r2 = 0
			}
		}

		if r1 != 0 {
			r.lineBuffer.WriteRune(r1)
		}
	}
}

func (r *Reader) readRune() (rune, error) {
	r1, _, err := r.input.ReadRune()

	// skip  BOM in the beginning of the file https://en.wikipedia.org/wiki/Byte_order_mark
	if r.rowCount == 1 && r.colCount == 0 && r1 == 0xFEFF {
		r1, _, err = r.input.ReadRune()
	}

	// consider \n, \r, \n\r, \r\n as just \n
	if r1 == '\r' || r1 == '\n' {
		r2, _, err := r.input.ReadRune()
		if err == nil && (r1 != '\n' || r2 != '\r') && (r1 != '\r' || r2 != '\n') {
			err = r.input.UnreadRune()
		}
		r1 = '\n'
	}

	if r1 == '\n' {
		r.rowCount++
		r.colCount = 0
	} else {
		r.colCount++
	}

	return r1, err
}
