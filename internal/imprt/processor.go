package imprt

import (
	"context"

	"highlite-parser/internal/csv"
	"highlite-parser/internal/highlite"
	"highlite-parser/internal/log"
)

// NewProcessor creates an Processor instance.
func NewProcessor(logger log.ILogger, productImport *ProductImport, client *highlite.Client) *Processor {
	return &Processor{
		logger:        logger,
		productImport: productImport,
		highClient:    client,
	}
}

// Processor handles highlite product update.
type Processor struct {
	logger        log.ILogger
	productImport *ProductImport
	highClient    *highlite.Client
}

// Update starts the update process.
func (p *Processor) Update(ctx context.Context) {
	p.logger.Debug("Getting items from highlite server")

	items, err := p.highClient.GetItemsReader(ctx)
	if err != nil {
		p.logger.Errorf("Can't get highlite items reader: %s", err.Error())

		return
	}

	csvParser := csv.NewReader(items)

	csvParser.ReadTitles()
	csvMapper := csv.NewTitleMap(csvParser.Titles())

	p.logger.Debug("Items processing start")

	i := 3 // temporary limit
	for run := true; run && i > 0; i-- {
		select {
		case <-ctx.Done():
			p.logger.Info("Context timeout")
			run = false

		default:
			if !csvParser.Next() {
				run = false
				break
			}

			pr := highlite.GetProductFromCSVImport(csvMapper, csvParser.Values())
			p.logger.Debugf("Processing product: category %s", pr.Category3.GetURL())

			if err := p.productImport.Import(ctx, pr); err != nil {
				p.logger.Errorf("Product processing error: %s", err.Error())
			}
		}
	}

	if csvParser.Err() != nil {
		p.logger.Errorf("Csv processing error: %s", csvParser.Err().Error())
	}

	p.logger.Debug("Stop items processing")
}
