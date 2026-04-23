package core_logger

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"time"
)

type Logger struct {
	*slog.Logger
	file *os.File
}

func NewLogger(config Config) (Logger, error) {
	if err := os.MkdirAll(config.Path, 0755); err != nil {
		return Logger{}, fmt.Errorf("failed to create logger folder: %w", err)
	}

	fileName := time.Now().UTC().Format("2006-01-02T05:04:15")
	fileName = fmt.Sprintf("%s/%s", config.Path, fileName)
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return Logger{}, fmt.Errorf("failed to create logger file: %w", err)
	}
	writter := io.MultiWriter(file, os.Stdout)

	opt := slog.HandlerOptions{
		Level: slog.LevelDebug,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				a.Value = slog.StringValue(time.Now().UTC().Format("2006-01-02T05:04:15.000000"))
			}
			return a
		},
	}

	handler := slog.NewTextHandler(writter, &opt)
	log := slog.New(handler)

	return Logger{
		Logger: log,
		file:   file,
	}, nil
}

func (l *Logger) Close() error {
	return l.file.Close()
}

func (l *Logger) Error(err error, args ...any) {
	l.Logger.Error(err.Error(), args...)
}

func GetLoggerFromContext(ctx context.Context) Logger {
	log := ctx.Value("logger")
	logger, ok := log.(Logger)
	if !ok {
		panic("failed to get logger from context")
	}
	return logger
}
