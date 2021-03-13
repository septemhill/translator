package logger

import (
	"context"
	"os"
)

type Logger interface {
	LogFields(context.Context, Fields)
	Println(context.Context, ...interface{})
	Printf(context.Context, string, ...interface{})
	Errorln(context.Context, ...interface{})
	Errorf(context.Context, string, ...interface{})
	Fatalln(context.Context, ...interface{})
	Fatalf(context.Context, string, ...interface{})
}

type contextLogger struct{}

func NewContext(ctx context.Context, logger Logger) context.Context {
	return context.WithValue(ctx, contextLogger{}, logger)
}

func LoggerFromContext(ctx context.Context) (Logger, bool) {
	logger, ok := ctx.Value(contextLogger{}).(Logger)
	return logger, ok
}

func ContextPrintln(ctx context.Context, a ...interface{}) {
	logger, ok := LoggerFromContext(ctx)
	if !ok {
		return
	}
	logger.Println(ctx, a...)
}

func ContextPrintf(ctx context.Context, format string, a ...interface{}) {
	logger, ok := LoggerFromContext(ctx)
	if !ok {
		return
	}
	logger.Printf(ctx, format, a...)
}

func ContextErrorln(ctx context.Context, a ...interface{}) {
	logger, ok := LoggerFromContext(ctx)
	if !ok {
		return
	}
	logger.Errorln(ctx, a...)
}

func ContextErrorf(ctx context.Context, format string, a ...interface{}) {
	logger, ok := LoggerFromContext(ctx)
	if !ok {
		return
	}
	logger.Errorf(ctx, format, a...)
}

func ContextFatalln(ctx context.Context, a ...interface{}) {
	logger, ok := LoggerFromContext(ctx)
	if !ok {
		return
	}
	logger.Fatalln(ctx, a...)
	os.Exit(1)
}

func ContextFatalf(ctx context.Context, format string, a ...interface{}) {
	logger, ok := LoggerFromContext(ctx)
	if !ok {
		return
	}
	logger.Fatalf(ctx, format, a...)
	os.Exit(1)
}
