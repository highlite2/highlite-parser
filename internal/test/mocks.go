package test

import (
	"highlite2-import/internal/cache"
	"highlite2-import/internal/highlite/image"
	"highlite2-import/internal/highlite/translation"
	"highlite2-import/internal/imprt"
	"highlite2-import/internal/log"
	"highlite2-import/internal/sylius"
)

// Logger ...
//go:generate mockery -name Logger
type Logger interface {
	log.ILogger
}

// ProductImport ...
//go:generate mockery -name ProductImport
type ProductImport interface {
	imprt.IProductImport
}

// SyliusClient ...
//go:generate mockery -name SyliusClient
type SyliusClient interface {
	sylius.IClient
}

// CacheMemo ...
//go:generate mockery -name CacheMemo
type CacheMemo interface {
	cache.IMemo
}

// TranslationDictionay ...
//go:generate mockery -name TranslationDictionay
type TranslationDictionay interface {
	translation.IDictionary
}

// ImageProvider ...
//go:generate mockery -name ImageProvider
type ImageProvider interface {
	image.IProvider
}

// CategoryImport ...
//go:generate mockery -name CategoryImport
type CategoryImport interface {
	imprt.ICategoryImport
}
