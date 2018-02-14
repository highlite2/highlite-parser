package test

import (
	"highlite-parser/internal/imprt"
	"highlite-parser/internal/log"
)

// Logger ...
//go:generate mockery -name Logger -output mocks
type Logger interface {
	log.ILogger
}

// ProductImport ...
//go:generate mockery -name ProductImport -output mocks
type ProductImport interface {
	imprt.IProductImport
}
