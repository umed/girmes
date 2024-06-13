package logging

import (
	"context"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/umed/girmes/internal/util"
)

func NewLogger(level string) *zerolog.Logger {
	logLevel := util.Must(zerolog.ParseLevel(level))

	logger := zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.DateTime}).
		With().
		Timestamp().
		Caller().
		Logger()
	logger.Level(logLevel)
	return &logger
}

func WithLogger(ctx context.Context, logger *zerolog.Logger) context.Context {
	return logger.WithContext(ctx)
}

func L(ctx context.Context) *zerolog.Logger {
	return zerolog.Ctx(ctx)
}
