package repository

import (
	"context"

	"github.com/septemhill/translator/service"
)

type DictionaryRepository interface {
	Record(ctx context.Context, word, src, dst string, tr *service.Translation) error
	Query(ctx context.Context, word, src, dst string) (*service.Translation, error)
	Close(ctx context.Context) error
}
