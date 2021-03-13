package main

import (
	"context"

	"github.com/septemhill/translator/service"
)

type Dictionary interface {
	Lookup(ctx context.Context, word, src, dst string, example bool) (*service.Translation, error)
}

type dictionary struct {
	ws service.WordTranslateService
}

func (d *dictionary) Lookup(ctx context.Context, word, src, dst string, example bool) (*service.Translation, error) {
	tr, err := d.ws.WordLookup(ctx, word, src, dst, example)
	if err != nil {
		return nil, err
	}
	return tr, nil
}

var _ Dictionary = (*dictionary)(nil)

func NewDictionary(ws service.WordTranslateService) *dictionary {
	return &dictionary{ws: ws}
}
