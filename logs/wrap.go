package logs

import (
	"fmt"
	"github.com/inconshreveable/log15"
	"runtime"
)

type Logger struct {
	log log15.Logger
}

func NewLogger() *Logger {
	return &Logger{log: log15.New()}
}

func (l *Logger) Info(msg string, ctx ...interface{}) {
	if _, file, line, ok := runtime.Caller(1); ok {
		ctx = append(ctx, "caller", fmt.Sprintf("%s: %d", file, line))
	}
	l.log.Info(msg, ctx...)
}

func a() {
	logger := NewLogger()
	logger.Info("hello", "a", "b")
	logger.Info("world")
}

func b() {
	fmt.Println("hello")
	a()
}
