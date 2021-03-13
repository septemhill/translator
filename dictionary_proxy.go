package main

import (
	"context"

	"github.com/septemhill/translator/logger"
	"github.com/septemhill/translator/repository"
	"github.com/septemhill/translator/service"
)

type dictionaryProxy struct {
	d    Dictionary
	repo repository.DictionaryRepository
}

func NewDictionaryProxy(d Dictionary, repo repository.DictionaryRepository) *dictionaryProxy {
	return &dictionaryProxy{
		d:    d,
		repo: repo,
	}
}

func (p *dictionaryProxy) Lookup(ctx context.Context, word, src, dst string, example bool) (*service.Translation, error) {
	dtr, err := p.repo.Query(ctx, word, src, dst)
	if err != nil {
		return nil, err
	}

	if dtr != nil {
		return dtr, nil
	}

	htr, err := p.d.Lookup(ctx, word, src, dst, example)
	if err != nil {
		return nil, err
	}
	if err := p.repo.Record(ctx, word, src, dst, htr); err != nil {
		logger.ContextErrorln(ctx, "Failed to record translation in local db: ", err)
	}
	return htr, nil
}

var _ Dictionary = (*dictionaryProxy)(nil)
