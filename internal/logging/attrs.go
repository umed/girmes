package logging

import "log/slog"

func Err(err error) slog.Attr {
	return slog.String("error", err.Error())
}
