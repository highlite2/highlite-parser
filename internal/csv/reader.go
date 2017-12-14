package csv

import (
	"bufio"
	"io"
	"strings"

	"github.com/pkg/errors"
)

var (
	errNoTitles         = errors.New("failed to get titles")
	errWrongColumnCount = errors.New("column count is wrong")
)

// NewReader creates a new Reader.
func NewReader(reader io.Reader) *Reader {
	return &Reader{
		scanner:   bufio.NewScanner(reader),
		separator: ";",
	}
}

// Reader reads highlite csv file.
type Reader struct {
	scanner   *bufio.Scanner
	separator string

	values []string
	titles []string

	err error
}

// SetSeparator sets separator to split csv row. Default value is ";".
func (c *Reader) SetSeparator(s string) {
	c.separator = s
}

// Err returns an error that might has occurred.
func (c *Reader) Err() error {
	return c.err
}

// Values returns last read line values.
func (c *Reader) Values() []string {
	return c.values
}

// Titles returns titles.
func (c *Reader) Titles() []string {
	return c.titles
}

// ReadTitles reads values from current line and saves them as titles.
// It returns true if the line was successfully read. False is returned
// when an error occurred on file end was reached.
func (c *Reader) ReadTitles() bool {
	if !c.scanner.Scan() {
		return c.handleErr(c.scanner.Err())
	}

	if c.titles = c.parseCSVLine(c.scanner.Text()); len(c.titles) == 0 {
		return c.handleErr(errNoTitles)
	}

	return true
}

// Next triggers next line reading. It returns true, if the line was successfully read.
// False is returned when an error occurred or file end was reached.
func (c *Reader) Next() bool {
	if !c.scanner.Scan() {
		return c.handleErr(c.scanner.Err())
	}

	c.values = c.parseCSVLine(c.scanner.Text())

	if len(c.values) == 0 || (len(c.titles) > 0 && len(c.titles) != len(c.values)) {
		return c.handleErr(errWrongColumnCount)
	}

	return true
}

// Saves the error and returns false.
func (c *Reader) handleErr(err error) bool {
	c.err = err

	return false
}

// Splits string into slice
func (c *Reader) parseCSVLine(str string) []string {
	return strings.Split(str, c.separator)
}
