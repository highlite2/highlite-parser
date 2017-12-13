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

// GetWindows1257Decoder decorates reader with Windows1257 decoder transformation.
func GetWindows1257Decoder(reader io.Reader) io.Reader {
	return transform.NewReader(reader, charmap.Windows1257.NewDecoder())
}

// NewCSVReader creates a new CSVReader.
func NewCSVReader(reader io.Reader, log internal.ILogger) *CSVReader {
	return &CSVReader{
		scanner: bufio.NewScanner(reader),
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

// Err returns an error that might has occurred.
func (c *CSVReader) Err() error {
	return c.err
}

// Values returns last read line values.
func (c *CSVReader) Values() []string {
	return c.values
}

// Titles returns titles.
func (c *CSVReader) Titles() []string {
	return c.titles
}

// ReadTitles reads values from current line and saves them as titles.
// It returns true is the line was successfully read. False is returned
// when an error occurred on file end was reached.
func (c *CSVReader) ReadTitles() bool {
	if !c.scanner.Scan() {
		return c.handleErr(c.scanner.Err())
	}

	if c.titles = parseCSVLine(c.scanner.Text()); len(c.titles) == 0 {
		return c.handleErr(errNoTitles)
	}

	return true
}

// Next triggers next line reading. It returns true, if the line was successfully read.
// False is returned when an error occurred on file end was reached.
func (c *CSVReader) Next() bool {
	if !c.scanner.Scan() {
		return c.handleErr(c.scanner.Err())
	}

	c.values = parseCSVLine(c.scanner.Text())

	if len(c.values) == 0 || (len(c.titles) > 0 && len(c.titles) != len(c.values)) {
		return c.handleErr(errWrongColumnCount)
	}

	return true
}

// Saves the error and returns false.
func (c *CSVReader) handleErr(err error) bool {
	c.err = err

	return false
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
