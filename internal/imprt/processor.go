package imprt

import (
	"context"
	"io"

	"highlite2-import/internal/csv"
	"highlite2-import/internal/highlite"
	"highlite2-import/internal/log"
	"highlite2-import/internal/queue"
)

// CSVParser is a highlite file updates parser interface
type CSVParser interface {
	GetNext() []string
	Next() bool
	Err() error
	Values() []string
	CurrentRowIndex() int
}

// NewProcessor creates an Processor instance.
func NewProcessor(logger log.ILogger, pool queue.IPool, productImport IProductImport, csvParser CSVParser) *Processor {
	return &Processor{
		logger:        logger,
		workerPool:    pool,
		productImport: productImport,
		titles:        highlite.GetRequiredCSVTitles(),
		csvParser:     csvParser,
	}
}

// Processor handles highlite product update.
type Processor struct {
	logger        log.ILogger
	workerPool    queue.IPool
	productImport IProductImport
	items         io.Reader
	titles        []string
	csvParser     CSVParser
}

// SetTitles sets titles
func (p *Processor) SetTitles(titles []string) {
	p.titles = titles
}

// Update starts the update process.
func (p *Processor) Update(ctx context.Context) {
	p.logger.Debug("starting update")

	csvMapper := csv.NewTitleMap(p.csvParser.GetNext())

	if err := csvMapper.CheckTitles(p.titles); err != nil {
		p.logger.Error(err.Error())
		return
	}

	i := 0
	for {
		select {
		case <-ctx.Done():
			p.logger.Warn("context timeout")
			return

		default:
			if !p.csvParser.Next() {
				if p.csvParser.Err() != nil {
					p.logger.Errorf("error processing csv with product updates: %s", p.csvParser.Err())
				}
				return
			}

			product := highlite.GetProductFromCSVImport(csvMapper, p.csvParser.Values())
			if err := csvMapper.CheckValues(p.csvParser.Values()); err != nil {
				p.logger.Errorf(
					"error processing csv with product updates: check values error on line %d: %s",
					p.csvParser.CurrentRowIndex(),
					err.Error(),
				)
				p.logger.Warn(product.String())
			} else {
				<-p.workerPool.AddJob(p.getImportJob(ctx, product))
			}
		}

		i++
		if i%50 == 0 {
			p.logger.Infof("processed %d products", i)
		}
	}
}

// Creates import job.
func (p *Processor) getImportJob(ctx context.Context, high highlite.Product) queue.IJob {
	return queue.NewCallbackJob(func() error {
		err := p.productImport.Import(ctx, high)
		if err != nil {
			p.logger.Errorf("product %s processing error: %s", high.No, err.Error())
		}
		return err
	})
}
