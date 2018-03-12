package translation

import (
	"fmt"
	"os"
	"testing"

	"highlite2-import/internal/csv"

	"github.com/stretchr/testify/assert"
)

var csvTestDataFilePath = "./csv_test_data.csv"
var csvTestDataProductIDs = []string{
	"20100", "20508", "20510", "20521", "20526", "20536", "20571", "20600", "20601", "20605", "20606", "20611",
	"20613", "20615", "20619", "20620", "20621", "20626", "20628", "20639", "20658", "20681", "20726", "30101",
}

func TestItemCSV(t *testing.T) {
	// arrange
	item := ProductCSV{
		MainText:   "--cattext_main_rus--",
		SubHeading: "--cattext_subh_rus--",
		USP:        "--USP_rus--",
		TechSpec:   "--techspec_rus--",
	}

	expectedDescription := item.USP + "\n" + item.MainText + "\n\n" + item.TechSpec
	expectedShortDescription := item.SubHeading

	// act
	// assert
	assert.Equal(t, expectedDescription, item.GetDescription())
	assert.Equal(t, expectedShortDescription, item.GetShortDescription())
}

func TestFillMemoryDictionaryFromCSV(t *testing.T) {
	// arrange
	dic := NewMemoryDictionary()
	loc := "ru_RU"
	items, itemsErr := readCSVTestDataFile()

	// act
	FillMemoryDictionaryFromCSV(dic, loc, csvTestDataFilePath, GetRussianTranslationsCSVTitles())

	// assert
	assert.Nil(t, itemsErr)
	for id, item := range items {
		pr, ok := dic.Get(loc, id)
		assert.True(t, ok)

		if ok {
			assert.Equal(t, item.GetDescription(), pr.GetDescription())
			assert.Equal(t, item.GetShortDescription(), pr.GetShortDescription())

			if id == "20100" {
				assert.Equal(t, ProductCSV{
					MainText:   "--cattext_main_rus--",
					SubHeading: "--cattext_subh_rus--",
					USP:        "--USP_rus--",
					TechSpec:   "--techspec_rus--",
				}, item)
			}
		}
	}
}

func readCSVTestDataFile() (map[string]ProductCSV, error) {
	file, err := os.Open(csvTestDataFilePath)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	var records = make(map[string]ProductCSV, 0)

	reader := csv.NewReader(file)
	titles := csv.NewTitleMap(reader.GetNext())

	i := 0
	for reader.Next() {
		if i >= len(csvTestDataProductIDs) {
			return nil, fmt.Errorf("reading unexpected row %d", i+1)
		}

		values := reader.Values()
		productNo := titles.GetString("Ordercode", values)
		if csvTestDataProductIDs[i] != productNo {
			return nil, fmt.Errorf("unexpected product No %s, expected %s", productNo, csvTestDataProductIDs[i])
		}

		records[productNo] = ProductCSV{
			MainText:   titles.GetString("cattext_main_rus", values),
			SubHeading: titles.GetString("cattext_subh_rus", values),
			USP:        titles.GetString("USP_rus", values),
			TechSpec:   titles.GetString("techspec_rus", values),
		}

		i++
	}

	if len(csvTestDataProductIDs) != len(records) {
		return nil, fmt.Errorf("unexpected product count %d, expected %d", len(records), len(csvTestDataProductIDs))
	}

	return records, reader.Err()
}
