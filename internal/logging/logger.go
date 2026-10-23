package logging

import (
	"io"
	"log/slog"
)

const millisecondTimeFormat = "2006-01-02T15:04:05.000Z07:00"

func New(out io.Writer) *slog.Logger {
	handler := slog.NewJSONHandler(out, &slog.HandlerOptions{
		ReplaceAttr: func(groups []string, attr slog.Attr) slog.Attr {
			if len(groups) == 0 && attr.Key == slog.TimeKey {
				return slog.String(
					slog.TimeKey,
					attr.Value.Time().Format(millisecondTimeFormat),
				)
			}

			return attr
		},
	})

	return slog.New(handler)
}
