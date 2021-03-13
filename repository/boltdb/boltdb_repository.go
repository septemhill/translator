package boltdb

import (
	"context"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/boltdb/bolt"
	"github.com/septemhill/translator/decorator"
	"github.com/septemhill/translator/decorator/compress"
	"github.com/septemhill/translator/decorator/serialization"
	"github.com/septemhill/translator/repository"
	"github.com/septemhill/translator/service"
)

var (
	tslBaseLoc = filepath.Join(os.Getenv("HOME"), ".tsl")
	tslDbLoc   = filepath.Join(tslBaseLoc, "db")
	// tslConfigLoc = filepath.Join(tslBaseLoc, "tsl_config.yaml")
)

type boltDbRepository struct {
	db *bolt.DB
}

func (repo *boltDbRepository) Record(ctx context.Context, word, src, dst string, tr *service.Translation) error {
	fn := func(tx *bolt.Tx) error {
		bk, err := tx.CreateBucketIfNotExists([]byte(word))
		if err != nil {
			return err
		}

		r := compress.Compress(serialization.Serialization(tr)).Decorate()
		b, _ := ioutil.ReadAll(r)
		return bk.Put([]byte(src+"-"+dst), b)
	}

	return repo.db.Update(fn)
}

func (repo *boltDbRepository) Query(ctx context.Context, word, src, dst string) (*service.Translation, error) {
	var tr *service.Translation
	var t service.Translation

	if err := repo.db.View(func(tx *bolt.Tx) error {
		bk := tx.Bucket([]byte(word))
		if bk == nil {
			return nil
		}

		b := bk.Get([]byte(src + "-" + dst))
		if len(b) == 0 {
			return nil
		}

		if err := serialization.Deserialization(compress.Decompress(decorator.NewBufferDecorator(b))).Decode(&t); err != nil {
			return err
		}
		tr = &t
		return nil
	}); err != nil {
		return nil, err
	}

	return tr, nil
}

func (repo *boltDbRepository) Close(ctx context.Context) error {
	return repo.db.Close()
}

func NewBoltDBRepository() (*boltDbRepository, error) {
	db, err := bolt.Open(tslDbLoc, 0600, nil)
	if err != nil {
		return nil, err
	}

	return &boltDbRepository{db: db}, nil
}

var _ repository.DictionaryRepository = (*boltDbRepository)(nil)
