package imprt

import (
	"sync"

	"strings"

	"highlite2-import/internal/cache"
	"highlite2-import/internal/highlite"
	"highlite2-import/internal/log"
	"highlite2-import/internal/sylius"
)

// IAttributesImport is a Sylius attributes importer.
type IAttributesImport interface {
	GetBrandID(high highlite.Product) (string, bool)
	//AddBrand(ctx context.Context, high highlite.Product) error
}

// NewAttributesImport creates new Sylius attributes importer.
func NewAttributesImport(client sylius.IClient, memo cache.IMemo, logger log.ILogger) *AttributesImport {
	return &AttributesImport{
		logger:   logger,
		client:   client,
		memo:     memo,
		lock:     sync.RWMutex{},
		brandMap: make(map[string]string),
	}
}

// AttributesImport is a Sylius attributes importer implementation.
type AttributesImport struct {
	logger   log.ILogger
	client   sylius.IClient
	memo     cache.IMemo
	lock     sync.RWMutex
	brandMap map[string]string
}

// GetBrandID gets brand ID
func (i *AttributesImport) GetBrandID(high highlite.Product) (string, bool) {
	i.lock.RLock()
	defer i.lock.RUnlock()

	brand := strings.ToLower(high.Brand)
	if id, ok := i.brandMap[brand]; ok {
		return id, true
	}

	return "", false
}
