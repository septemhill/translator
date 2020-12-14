package nodb

import (
	"context"

	"github.com/septemhill/translator/repository"
	"github.com/septemhill/translator/service"
)

type noDbRepository struct{}

func (repo *noDbRepository) Record(ctx context.Context, word, src, dst string, tr *service.Translation) error {
	return nil
}

func (repo *noDbRepository) Query(ctx context.Context, word, src, dst string) (*service.Translation, error) {
	return nil, nil
}

func (repo *noDbRepository) Close(ctx context.Context) error {
	return nil
}

func NewNoDBRepository() *noDbRepository {
	return &noDbRepository{}
}

var _ repository.DictionaryRepository = (*noDbRepository)(nil)
