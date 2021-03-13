package filelog

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/septemhill/translator/logger"
)

const (
	logTag     = "[LOG]"
	infoTag    = "[INFO]"
	errorTag   = "[ERROR]"
	fatalTag   = "[FATAL]"
	timeFormat = "2006-01-02T15:04:05"
	logPath    = "."
)

var onceLogger sync.Once
var flogger *fileLogger

type fileLogger struct {
	mu sync.Mutex
	w  *bufio.Writer
}

func FileLogger() *fileLogger {
	onceLogger.Do(func() {
		if flogger == nil {
			flogger = newFileLogger()
		}
	})
	return flogger
}

func newFileLogger() *fileLogger {
	fname := fmt.Sprintf("%s/%s_%d_meet_server.log", logPath, time.Now().Format(timeFormat), os.Getpid())
	f, err := os.Create(fname)
	if err != nil {
		log.Fatal("Failed to create file logger: ", err)
	}
	w := bufio.NewWriter(f)
	return &fileLogger{w: w}
}

func (l *fileLogger) addLogPrefix(tag string) {
	fmt.Fprintf(l.w, "%s %s ", tag, time.Now().Format(timeFormat))
}

func (l *fileLogger) outputln(tag string, a ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.addLogPrefix(tag)
	fmt.Fprintln(l.w, a...)
	l.w.Flush()
}

func (l *fileLogger) outputf(tag string, format string, a ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.addLogPrefix(tag)
	fmt.Fprintf(l.w, format, a...)
	l.w.Flush()
}

func (l *fileLogger) LogFields(_ context.Context, f logger.Fields) {
	l.outputln(logTag)
}

func (l *fileLogger) Println(_ context.Context, a ...interface{}) {
	l.outputln(infoTag, a...)
}

func (l *fileLogger) Printf(_ context.Context, format string, a ...interface{}) {
	l.outputf(infoTag, format, a...)
}

func (l *fileLogger) Errorln(_ context.Context, a ...interface{}) {
	l.outputln(errorTag, a...)
}

func (l *fileLogger) Errorf(_ context.Context, format string, a ...interface{}) {
	l.outputf(errorTag, format, a...)
}

func (l *fileLogger) Fatalln(_ context.Context, a ...interface{}) {
	l.outputln(fatalTag, a...)
}

func (l *fileLogger) Fatalf(_ context.Context, format string, a ...interface{}) {
	l.outputf(fatalTag, format, a...)
}

var _ logger.Logger = (*fileLogger)(nil)
