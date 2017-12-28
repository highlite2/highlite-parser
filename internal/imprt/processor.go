package imprt

import (
	"context"

	"highlite-parser/internal/csv"
	"highlite-parser/internal/highlite"
	"highlite-parser/internal/log"
	"highlite-parser/internal/queue"
)

// NewProcessor creates an Processor instance.
func NewProcessor(logger log.ILogger, pool *queue.Pool, productImport *ProductImport, client *highlite.Client) *Processor {
	return &Processor{
		logger:        logger,
		workerPool:    pool,
		productImport: productImport,
		highClient:    client,
	}
}

// Processor handles highlite product update.
type Processor struct {
	logger        log.ILogger
	workerPool    *queue.Pool
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

	p.logger.Debug("CSV parsing start")

	for run := true; run; {
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

	p.logger.Debug("CSV parsing stop")
}
