package prelog

import (
	"context"
	"fmt"
	"sync"

	"github.com/septemhill/translator/logger"
)

var onceLogger sync.Once
var plogger *preLogger

type preLogger struct {
	mu sync.Mutex
	l  logger.Logger
	p  string
}

func PreLogger(l logger.Logger, prefix string) *preLogger {
	onceLogger.Do(func() {
		if plogger == nil {
			plogger = newPreLogger(l, prefix)
		}
	})
	return plogger
}

func newPreLogger(l logger.Logger, prefix string) *preLogger {
	return &preLogger{l: l, p: prefix}
}

func (l *preLogger) LogFields(ctx context.Context, f logger.Fields) {
	l.mu.Lock()
	defer l.mu.Unlock()
	n := append([]interface{}{l.p}, f)
	l.l.Println(ctx, n...)
}

func (l *preLogger) Println(ctx context.Context, a ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	n := append([]interface{}{l.p}, a...)
	l.l.Println(ctx, n...)
}

func (l *preLogger) Printf(ctx context.Context, format string, a ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.l.Printf(ctx, fmt.Sprintf("%s %s", l.p, format), a...)
}

func (l *preLogger) Errorln(ctx context.Context, a ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	n := append([]interface{}{l.p}, a...)
	l.l.Errorln(ctx, n...)
}

func (l *preLogger) Errorf(ctx context.Context, format string, a ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.l.Errorf(ctx, fmt.Sprintf("%s %s", l.p, format), a...)
}

func (l *preLogger) Fatalln(ctx context.Context, a ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	n := append([]interface{}{l.p}, a...)
	l.l.Fatalln(ctx, n...)
}

func (l *preLogger) Fatalf(ctx context.Context, format string, a ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.l.Fatalf(ctx, fmt.Sprintf("%s %s", l.p, format), a...)
}

var _ logger.Logger = (*preLogger)(nil)
