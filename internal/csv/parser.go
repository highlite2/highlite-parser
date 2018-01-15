package csv

import (
	"bufio"
	"io"
)

type Parser struct {
	Separator rune

	r            *bufio.Reader
	parsingError error
}

func NewParser(r io.Reader) *Parser {
	return &Parser{
		Separator: ',',

		r: bufio.NewReader(r),
	}
}

func (p *Parser) readRune() (rune, error) {
	r1, _, err := p.r.ReadRune()

	if r1 == '\r' {
		r1, _, err = p.r.ReadRune()
		if err == nil {
			if r1 != '\n' {
				p.r.UnreadRune()
				r1 = '\r'
			}
		}
	}

	return r1, err
}

func (p *Parser) Next() ([]string, error) {
	if p.parsingError != nil {
		return nil, p.parsingError
	}

	fields, err := p.parseRecord()
	if err != nil {
		p.parsingError = err
		if err != io.EOF || len(fields) == 0 {
			return nil, err
		}
	}

	return fields, nil
}

func (p *Parser) parseRecord() ([]string, error) {
	fields := make([]string, 0)
	for {
		ru, err := p.readRune()
		if err != nil {
			if len(fields) > 0 {
				fields = append(fields, "")
			}

			return fields, err
		}

		switch ru {
		case '\n':
			return append(fields, ""), nil

		case p.Separator:
			fields = append(fields, "")

		default:
			field, stop, err := p.getField(ru)
			if field != nil {
				fields = append(fields, *field)
			}
			if err != nil {
				return fields, err
			}
			if stop {
				return fields, nil
			}
		}
	}
}

func (p *Parser) getField(ru rune) (*string, bool, error) {
	switch ru {
	case '"', '\'':
		return p.getQuotedField(ru)
	default:
		return p.getUnquotedField(ru)
	}
}

func (p *Parser) getUnquotedField(r1 rune) (*string, bool, error) {
	field := string(r1)
	for {
		ru, err := p.readRune()
		if err != nil {
			return &field, true, err
		}

		switch ru {
		case '\n':
			return &field, true, nil
		case p.Separator:
			return &field, false, nil
		}

		field += string(ru)
	}
}

func (p *Parser) getQuotedField(pair rune) (*string, bool, error) {
	field := ""

	for {
		ru, err := p.readRune()
		if err != nil {
			return &field, true, err
		}

		if ru == pair {
			r2, err := p.readRune()
			if err != nil {
				return &field, true, err
			}

			switch r2 {
			case '\n':
				return &field, true, nil
			case p.Separator:
				return &field, false, nil
			case pair:
				// it was a quoted quote
			default:
				if err := p.r.UnreadRune(); err != nil {
					return &field, true, err
				}
			}
		}

		field += string(ru)
	}
}
