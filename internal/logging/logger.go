package logging

import (
	"context"
	"log/slog"
	"os"

	"github.com/umed/girmes/internal/util"
)

const LevelFatal = slog.LevelError + 1

type CustomLogger struct {
	*slog.Logger
}

var DefaultLogger = NewLogger("debug")

func (l *CustomLogger) Fatal(msg string, args ...any) {
	l.Log(context.Background(), LevelFatal, msg, args...)
	os.Exit(1)
}

func ParseLevel(level string) (slog.Level, error) {
	var logLevel slog.Level
	err := logLevel.UnmarshalText([]byte(level))
	return logLevel, err
}

func MustParseLevel(level string) slog.Level {
	return util.Must(ParseLevel(level))
}

func customReplaceAttr(groups []string, a slog.Attr) slog.Attr {
	if a.Key == slog.LevelKey && a.Value.Any().(slog.Level) == LevelFatal {
		a.Value = slog.StringValue("FATAL")
	}
	return a
}

func NewLogger(level string) *CustomLogger {
	return &CustomLogger{
		Logger: slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level:       MustParseLevel(level),
			ReplaceAttr: customReplaceAttr,
		})),
	}
}

type ctxKey struct{}

func Ctx(ctx context.Context, logger *CustomLogger) context.Context {
	return context.WithValue(ctx, ctxKey{}, &logger)
}

func L(ctx context.Context) *CustomLogger {
	if logger, ok := ctx.Value(ctxKey{}).(*CustomLogger); ok {
		return logger
	}
	return DefaultLogger
}
