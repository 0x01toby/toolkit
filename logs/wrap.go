package logs

import (
	"context"
	"fmt"
	"github.com/inconshreveable/log15"
	"go.opentelemetry.io/otel/trace"
	"reflect"
	"runtime"
	"strings"
)

const (
	TraceID  = "trace_id"
	SpanID   = "span_id"
	NoopSpan = "noopSpan"
	Caller   = "caller"
)

type Logger struct {
	log log15.Logger
}

func NewLogger(k, v string) *Logger {
	return &Logger{log: log15.New(k, v)}
}

func generateKvPair(ctx context.Context, kvPair ...interface{}) []interface{} {
	elems := make([]interface{}, 0)
	span := trace.SpanFromContext(ctx)
	if !strings.EqualFold(reflect.TypeOf(span).Name(), NoopSpan) {
		elems = append(elems, TraceID, span.SpanContext().TraceID().String(), SpanID, span.SpanContext().SpanID().String())
	}
	elems = append(elems, kvPair...)
	if _, file, line, ok := runtime.Caller(2); ok {
		elems = append(elems, Caller, fmt.Sprintf("%s: %d", file, line))
	}
	return elems
}

func (l *Logger) Debug(ctx context.Context, msg string, kvPair ...interface{}) {
	l.log.Debug(msg, generateKvPair(ctx, kvPair...)...)
}

func (l *Logger) Info(ctx context.Context, msg string, kvPair ...interface{}) {
	l.log.Info(msg, generateKvPair(ctx, kvPair...)...)
}

func (l *Logger) Warn(ctx context.Context, msg string, kvPair ...interface{}) {
	l.log.Warn(msg, generateKvPair(ctx, kvPair...)...)
}

func (l *Logger) Error(ctx context.Context, msg string, kvPair ...interface{}) {
	l.log.Error(msg, generateKvPair(ctx, kvPair...)...)
}
