package stdoutlog

import (
	"context"
	"encoding/json"
	"log"
	"sync"

	"github.com/septemhill/translator/logger"
)

var onceLogger sync.Once
var slogger *stdoutLogger

type stdoutLogger struct {
	mu sync.Mutex
}

func StdoutLogger() *stdoutLogger {
	onceLogger.Do(func() {
		if slogger == nil {
			slogger = newStdoutLogger()
		}
	})
	return slogger
}

func newStdoutLogger() *stdoutLogger {
	return &stdoutLogger{}
}

func (l *stdoutLogger) LogFields(_ context.Context, f logger.Fields) {
	l.mu.Lock()
	defer l.mu.Unlock()
	d, _ := json.MarshalIndent(f, "", " ")
	log.Println(string(d))
}

func (l *stdoutLogger) Println(_ context.Context, a ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	log.Println(a...)
}

func (l *stdoutLogger) Printf(_ context.Context, format string, a ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	log.Printf(format, a...)
}

func (l *stdoutLogger) Errorln(_ context.Context, a ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	log.Println(a...)
}

func (l *stdoutLogger) Errorf(_ context.Context, format string, a ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	log.Printf(format, a...)
}

func (l *stdoutLogger) Fatalln(_ context.Context, a ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	log.Fatal(a...)
}

func (l *stdoutLogger) Fatalf(_ context.Context, format string, a ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	log.Fatalf(format, a...)
}

var _ logger.Logger = (*stdoutLogger)(nil)
