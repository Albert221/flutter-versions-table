package log

import (
	stdLog "log"
	"strings"
)

type Logger struct {
	depth int
}

func New() *Logger {
	return &Logger{depth: 0}
}

func (l *Logger) Sub() *Logger {
	return &Logger{depth: l.depth + 1}
}

func (l *Logger) Info(log string, args ...any) {
	prefix := strings.Repeat("  ", l.depth) + "• "
	stdLog.Printf(prefix+log, args...)
}
