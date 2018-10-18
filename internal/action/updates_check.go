package action

import (
	"context"
	"time"

	"highlite2-import/internal"
	"highlite2-import/internal/csv"
	"highlite2-import/internal/highlite"
	"highlite2-import/internal/log"
)

// UpdatesCheck checks product updates file
func UpdatesCheck(ctx context.Context, config internal.Config, logger log.ILogger) {
	ctx, cancel := context.WithTimeout(ctx, config.ImportTimeout)
	defer cancel()
	defer timeTrack(logger, time.Now(), "Import")

	csvParser, err := getHighliteProductUpdatesCSVParser(ctx, config, logger)
	if err != nil {
		logger.Errorf("cant get highlite items reader: %s", err)
		return
	}

	csvMapper := csv.NewTitleMap(csvParser.GetNext())
	if err := csvMapper.CheckTitles(highlite.GetRequiredCSVTitles()); err != nil {
		logger.Errorf(err.Error())
		return
	}

	var productCounter int

	for csvParser.Next() {
		productCounter++
		product := highlite.GetProductFromCSVImport(csvMapper, csvParser.Values())
		if err := csvMapper.CheckValues(csvParser.Values()); err != nil {
			logger.Errorf("check values on line %d: %s", csvParser.CurrentRowIndex(), err)
			logger.Warn(product.String())
		}
	}

	if csvParser.Err() != nil {
		logger.Errorf("error processing csv with product updates: %s", csvParser.Err())
	}

	logger.Infof("processed %d products", productCounter)
}
