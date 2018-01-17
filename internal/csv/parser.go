package csv

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
)

type Parser struct {
	Separator    rune
	QuotedQuotes bool
	FieldsCount  int
	FieldsFixed  bool

	input        *bufio.Reader
	parsingErr   error
	lineBuffer   bytes.Buffer
	fieldIndexes []int

	rowCount        int
	colCount        int
	currentRowIndex int
}

func NewParser(r io.Reader) *Parser {
	return &Parser{
		Separator:   ',',
		FieldsFixed: true,

		input:    bufio.NewReader(r),
		rowCount: 1,
	}
}

func (p *Parser) Next() ([]string, error) {
	if p.parsingErr != nil {
		return nil, p.parsingErr
	}

	if err := p.parseRecord(); err != nil {
		p.parsingErr = err
		if err != io.EOF || len(p.fieldIndexes) == 0 {
			return nil, err
		}
	}

	line := p.lineBuffer.String()
	fieldCount := len(p.fieldIndexes)
	fields := make([]string, fieldCount)

	for i, idx := range p.fieldIndexes {
		if i == fieldCount-1 {
			fields[i] = line[idx:]
		} else {
			fields[i] = line[idx:p.fieldIndexes[i+1]]
		}
	}

	if p.FieldsCount == 0 {
		p.FieldsCount = fieldCount
	} else if p.FieldsFixed && p.FieldsCount != fieldCount {
		p.parsingErr = fmt.Errorf("wrong column count, row index %d", p.currentRowIndex)

		return fields, p.parsingErr
	}

	return fields, nil
}

func (p *Parser) parseRecord() error {
	p.lineBuffer.Reset()
	p.fieldIndexes = p.fieldIndexes[:0]
	p.currentRowIndex = p.rowCount

	for {
		r1, err := p.readRune()
		if err != nil {
			if len(p.fieldIndexes) > 0 {
				p.setNextFieldIndex()
			}

			return err
		}

		switch r1 {
		case '\n':
			p.setNextFieldIndex()

			return nil

		case p.Separator:
			p.setNextFieldIndex()

		default:
			stop, err := p.getField(r1)
			if err != nil {
				return err
			}
			if stop {
				return nil
			}
		}
	}
}

func (p *Parser) setNextFieldIndex() {
	p.fieldIndexes = append(p.fieldIndexes, p.lineBuffer.Len())
}

func (p *Parser) getField(r1 rune) (bool, error) {
	p.setNextFieldIndex()

	switch r1 {
	case '"', '\'':
		return p.getQuotedField(r1)
	default:
		return p.getUnquotedField(r1)
	}
}

func (p *Parser) getUnquotedField(r1 rune) (bool, error) {
	var err error
	p.lineBuffer.WriteRune(r1)
	for {
		r1, err = p.readRune()
		if err != nil {
			return true, err
		}

		if r1 == '\n' || r1 == p.Separator {
			return r1 == '\n', nil
		}

		p.lineBuffer.WriteRune(r1)
	}
}

func (p *Parser) getQuotedField(pair rune) (bool, error) {
	var r1, r2 rune
	var err error

	for {
		r1 = r2
		r2, err = p.readRune()
		if err != nil {
			return true, err
		}

		if r2 == '\n' || r2 == p.Separator {
			if r1 == pair {
				return r2 == '\n', nil
			}
		}
		if p.QuotedQuotes {
			if r1 == pair && r2 != pair {
				return true, fmt.Errorf("unquoted quote on %d:%d", p.currentRowIndex, p.colCount-1)
			}

			if r1 == pair && r2 == pair {
				r2 = 0
			}
		}

		if r1 != 0 {
			p.lineBuffer.WriteRune(r1)
		}
	}
}

func (p *Parser) readRune() (rune, error) {
	r1, _, err := p.input.ReadRune()

	if r1 == '\r' {
		r1, _, err = p.input.ReadRune()
		if err == nil {
			if r1 != '\n' {
				p.input.UnreadRune()
				r1 = '\r'
			}
		}
	}

	if r1 == '\n' {
		p.rowCount++
		p.colCount = 0
	} else {
		p.colCount++
	}

	return r1, err
}
