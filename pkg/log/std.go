package log

import (
	"fmt"
	"log/slog"
)

type stdLogger struct {
	logger *slog.Logger
}

func NewStdLogger(logger *slog.Logger) *stdLogger {
	return &stdLogger{
		logger: logger,
	}
}

func (sl *stdLogger) Fatalf(format string, v ...any) {
	sl.logger.Error(fmt.Sprintf(format, v...))
}

func (sl *stdLogger) Printf(format string, v ...any) {
	sl.logger.Info(fmt.Sprintf(format, v...))
}
