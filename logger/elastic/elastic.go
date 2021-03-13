package elastic

import (
	"context"
	"log"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/septemhill/translator/logger"
)

type elasticLogger struct {
	es *elasticsearch.Client
}

func NewElasticSearchLog() *elasticLogger {
	es, err := elasticsearch.NewDefaultClient()
	if err != nil {
		log.Fatal("Failed to create elastic search client")
	}

	return &elasticLogger{
		es: es,
	}
}

func (l *elasticLogger) LogFields(context.Context, logger.Fields) {
}

func (l *elasticLogger) Println(context.Context, ...interface{}) {
}

func (l *elasticLogger) Printf(context.Context, string, ...interface{}) {
}

func (l *elasticLogger) Errorln(context.Context, ...interface{}) {
}

func (l *elasticLogger) Errorf(context.Context, string, ...interface{}) {
}

func (l *elasticLogger) Fatalln(context.Context, ...interface{}) {
}

func (l *elasticLogger) Fatalf(context.Context, string, ...interface{}) {
}

var _ logger.Logger = (*elasticLogger)(nil)
