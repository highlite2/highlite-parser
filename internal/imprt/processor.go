package imprt

import (
	"context"
	"io"

	"highlite-parser/internal/csv"
	"highlite-parser/internal/highlite"
	"highlite-parser/internal/log"
	"highlite-parser/internal/queue"
)

// NewProcessor creates an Processor instance.
func NewProcessor(logger log.ILogger, pool *queue.Pool, productImport *ProductImport, items io.Reader) *Processor {
	return &Processor{
		logger:        logger,
		workerPool:    pool,
		productImport: productImport,
		items:         items,
	}
}

// Processor handles highlite product update.
type Processor struct {
	logger        log.ILogger
	workerPool    *queue.Pool
	productImport *ProductImport
	items         io.Reader
}

// Update starts the update process.
func (p *Processor) Update(ctx context.Context) {
	p.logger.Debug("Starting update")

	csvParser := csv.NewReader(p.items)
	csvParser.ReadTitles()
	csvMapper := csv.NewTitleMap(csvParser.Titles())

	defer p.handleCSVParserFinish(csvParser)

	for {
		select {
		case <-ctx.Done():
			p.logger.Warn("Context timeout")
			return

		default:
			if !csvParser.Next() {
				return
			}

			product := highlite.GetProductFromCSVImport(csvMapper, csvParser.Values())
			<-p.workerPool.AddJob(p.getImportJob(ctx, product))
		}
	}
}

// Creates import job.
func (p *Processor) getImportJob(ctx context.Context, high highlite.Product) queue.IJob {
	return queue.NewCallbackJob(func() error {
		p.logger.Infof("Processing product: %s, category: %s", high.No, high.Category3.GetURL())

		err := p.productImport.Import(ctx, high)
		if err != nil {
			p.logger.Errorf("Product processing error: %s", err.Error())
		}

		return err
	})
}

// Handles csv parser finish.
func (p *Processor) handleCSVParserFinish(csvParser *csv.Reader) {
	if csvParser.Err() != nil {
		p.logger.Errorf("Csv processing error: %s", csvParser.Err().Error())
	}

	p.logger.Debug("Stop update")
}
