package imprt

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"strings"
	"sync"
	"testing"

	"highlite2-import/internal/highlite"
	"highlite2-import/internal/queue"
	"highlite2-import/internal/test/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var titles = []string{`Product No.`, `Product Name`, `Country of Origin`, `Unit Price`, `Category`}

func TestProcessor_Update(t *testing.T) {
	// arrange
	pool := queue.NewPool(10)
	logger := &mocks.Logger{}
	logger.On("Debug", mock.Anything)
	logger.On("Infof", mock.Anything)

	lock := sync.Mutex{}
	actual := make(map[string][]string)
	pimp := &mocks.ProductImport{}
	pimp.On("Import", mock.Anything, mock.Anything).Times(len(processorTestData)).Return(nil).Run(
		func(args mock.Arguments) {
			lock.Lock()
			defer lock.Unlock()
			p, _ := args.Get(1).(highlite.Product)
			actual[p.No] = []string{p.No, p.Name, p.Country, fmt.Sprintf("%.2f", p.Price), p.Category1.Name}
		},
	)

	// act
	processor := NewProcessor(logger, pool, pimp, processorTestDataGetReader())
	processor.SetTitles(titles)
	processor.Update(context.Background())
	<-pool.Stop()

	// assert
	assert.Equal(t, processorTestData, actual)
}

var processorTestData = map[string][]string{
	"1": {"1", "product 1", "country 1", "101.02", "fashion 1"},
	"2": {"2", "product 2", "country 2", "102.45", "fashion 2"},
	"3": {"3", "product 3", "country 3", "103.12", "fashion 3"},
	"4": {"4", "product 4", "country 4", "104.56", "fashion 4"},
	"5": {"5", "product 5", "country 5", "105.12", "fashion 5"},
	"6": {"6", "product 6", "country 6", "106.32", "fashion 6"},
	"7": {"7", "product 7", "country 7", "107.00", "fashion 7"},
	"8": {"8", "product 8", "country 8", "108.56", "fashion 8"},
}

func processorTestDataGetReader() io.Reader {
	writer := bytes.Buffer{}
	writer.WriteString(strings.Join(titles, ";"))
	writer.WriteByte('\n')

	for _, product := range processorTestData {
		writer.WriteString(strings.Join(product, ";"))
		writer.WriteRune('\n')
	}

	return bytes.NewReader(writer.Bytes())
}
