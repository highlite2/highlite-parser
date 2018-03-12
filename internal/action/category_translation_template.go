package action

import (
	"context"
	"fmt"
	"os"
	"strings"

	"highlite2-import/internal"
	"highlite2-import/internal/csv"
	"highlite2-import/internal/highlite"
	"highlite2-import/internal/log"

	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
)

type CategoryTranslationTemplate struct {
}

func (c *CategoryTranslationTemplate) Do(ctx context.Context, config internal.Config, logger log.ILogger) error {
	outputFileName := "categories.csv"

	file, err := os.Open(config.ItemsFilePath)
	if err != nil {
		return fmt.Errorf("can't open file for reading items: %s", err.Error())
	}
	defer file.Close()

	itemsReader := transform.NewReader(file, charmap.Windows1257.NewDecoder())
	csvParser := csv.NewReader(itemsReader)
	csvParser.Separator = ';'
	csvMapper := csv.NewTitleMap(csvParser.GetNext())

	output, err := os.Create(outputFileName)
	if err != nil {
		return fmt.Errorf("can't create file for writing %s", outputFileName)
	}
	defer output.Close()

	if err := c.writeTitles(output); err != nil {
		return err
	}

	readyMap := map[string]bool{}

	i := 0
	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("context timeout")

		default:
			if !csvParser.Next() {
				if csvParser.Err() != nil {
					return fmt.Errorf("Error processing csv with product updates: %s", csvParser.Err())
				}

				return nil
			}

			product := highlite.GetProductFromCSVImport(csvMapper, csvParser.Values())
			if err := c.writeCategory(output, readyMap, product.Category1); err != nil {
				return err
			}

			if err := c.writeCategory(output, readyMap, product.Category2); err != nil {
				return err
			}

			if err := c.writeCategory(output, readyMap, product.Category3); err != nil {
				return err
			}
		}

		i++
		if i%100 == 0 {
			logger.Infof("Processed %d products", i)
		}
	}
}

func (c *CategoryTranslationTemplate) writeCategory(output *os.File, cache map[string]bool, cat *highlite.Category) error {
	if _, exists := cache[cat.GetCode()]; exists {
		return nil
	}

	data := []string{
		cat.GetCode(),
		cat.GetURL(),
		cat.Name,
		"\n",
	}

	if _, err := output.WriteString(strings.Join(data, ";")); err != nil {
		return err
	}

	cache[cat.GetCode()] = true

	return nil
}

func (c *CategoryTranslationTemplate) writeTitles(output *os.File) error {
	data := []string{
		"Code",
		"URL",
		"en_US",
		"ru_RU\n",
	}

	if _, err := output.WriteString(strings.Join(data, ";")); err != nil {
		return err
	}

	return nil
}
