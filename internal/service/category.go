package service

import (
	"highlite-parser/internal/cache"
	"highlite-parser/internal/client/sylius"
	"unicode"
)


type category struct {
	NameEn string
	NameRu string
}

type categoryRepository struct {
	client sylius.IClient
	memo   cache.IMemo
}

func (c *categoryRepository) GetLeaf(categories ...category) {

}
