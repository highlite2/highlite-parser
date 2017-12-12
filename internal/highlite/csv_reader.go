package highlite

import (
	"bufio"
	"io"
	"strings"

	"highlite-parser/internal"

	"github.com/pkg/errors"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
)

var (
	errNoTitles         = errors.New("failed to get titles")
	errWrongColumnCount = errors.New("column count is wrong")
)

// NewCSVReader creates a new CSVReader.
func NewCSVReader(reader io.Reader, log internal.ILogger) *CSVReader {
	return &CSVReader{
		scanner: bufio.NewScanner(reader),
		log:     log,
	}
}

// NewCSVReaderWithWindows1257Decoder creates a new CSVReader and decorates
// reader with Windows1257 decoder
func NewCSVReaderWithWindows1257Decoder(reader io.Reader, log internal.ILogger) *CSVReader {
	return &CSVReader{
		scanner: bufio.NewScanner(transform.NewReader(reader, charmap.Windows1257.NewDecoder())),
		log:     log,
	}
}

// CSVReader reads highlite csv file.
type CSVReader struct {
	scanner *bufio.Scanner
	log     internal.ILogger

	values []string
	titles []string

	err error
}

// Err returns an error that might has occurred
func (c *CSVReader) Err() error {
	return c.err
}

// Values returns last row
func (c *CSVReader) Values() []string {
	return c.values
}

// Titles returns titles
func (c *CSVReader) Titles() []string {
	return c.titles
}

// Next triggers next line reading. It returns true, is a line was successfully read.
// False is returned when an error occurred on end of file was reached.
func (c *CSVReader) Next() bool {
	if err := c.extractTitles(); err != nil {
		return c.handleErr(err)
	}

	if !c.scanner.Scan() {
		return c.handleErr(c.scanner.Err())
	}

	c.values = parseCSVLine(c.scanner.Text())
	if len(c.titles) != len(c.values) {
		return c.handleErr(errWrongColumnCount)
	}

	return true
}

// Saved an error
func (c *CSVReader) handleErr(err error) bool {
	c.err = err

	return false
}

// Extracts titles
func (c *CSVReader) extractTitles() error {
	if len(c.titles) > 0 {
		return nil
	}

	if !c.scanner.Scan() {
		return c.scanner.Err()
	}

	if c.titles = parseCSVLine(c.scanner.Text()); len(c.titles) == 0 {
		return errNoTitles
	}

	return nil
}

// Splits string into slice and removes quotes
func parseCSVLine(str string) []string {
	slice := strings.Split(str, ";")
	for i, v := range slice {
		v = strings.Trim(v, `"`)
		slice[i] = v
	}

	return slice
}
