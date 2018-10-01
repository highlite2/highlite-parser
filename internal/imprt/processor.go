package imprt

import (
	"context"
	"io"

	"highlite2-import/internal/csv"
	"highlite2-import/internal/highlite"
	"highlite2-import/internal/log"
	"highlite2-import/internal/queue"
)

// NewProcessor creates an Processor instance.
func NewProcessor(logger log.ILogger, pool queue.IPool, productImport IProductImport, items io.Reader) *Processor {
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
	workerPool    queue.IPool
	productImport IProductImport
	items         io.Reader
}

// Update starts the update process.
func (p *Processor) Update(ctx context.Context) {
	p.logger.Debug("Starting update")

	csvParser := csv.NewReader(p.items)
	csvParser.Separator = ';'
	csvParser.FieldsFixed = false
	csvParser.OneRowRecord = true
	csvMapper := csv.NewTitleMap(csvParser.GetNext())

	if err := csvMapper.CheckTitles(highlite.GetRequiredCSVTitles()); err != nil {
		p.logger.Error(err.Error())
		return
	}

	i := 0
	for {
		select {
		case <-ctx.Done():
			p.logger.Warn("Context timeout")

			return

		default:
			if !csvParser.Next() {
				if csvParser.Err() != nil {
					p.logger.Errorf("Error processing csv with product updates: %s", csvParser.Err().Error())
				}

				return
			}

			product := highlite.GetProductFromCSVImport(csvMapper, csvParser.Values())
			if err := csvMapper.CheckValues(csvParser.Values()); err != nil {
				p.logger.Errorf(
					"Error processing csv with product updates: check values error on line %d: %s",
					csvParser.CurrentRowIndex(),
					err.Error(),
				)
				p.logger.Warn(product.String())
			} else {
				<-p.workerPool.AddJob(p.getImportJob(ctx, product))
			}
		}

		i++
		if i%50 == 0 {
			p.logger.Infof("Processed %d products", i)
		}
	}
}

// Creates import job.
func (p *Processor) getImportJob(ctx context.Context, high highlite.Product) queue.IJob {
	return queue.NewCallbackJob(func() error {
		err := p.productImport.Import(ctx, high)
		if err != nil {
			p.logger.Errorf("Product %s processing error: %s", high.No, err.Error())
		}

		return err
	})
}
