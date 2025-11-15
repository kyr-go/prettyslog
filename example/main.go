package main

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/kyr-go/prettyslog"
)

func main() {
	newHandler := prettyslog.NewHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.Level(-8),
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				t := a.Value.Any().(time.Time)
				formatted := t.Format(time.Kitchen)
				a.Value = slog.StringValue(formatted)
			}
			if a.Key == slog.LevelKey {
				if a.Value.Any().(slog.Level) == slog.Level(-8) {
					return slog.Any(slog.LevelKey, "EMER")
				}
			}
			return a
		},
	}).WithGroup("HandlerOuter").WithAttrs([]slog.Attr{slog.Group("HandlerInner", "key", "value")})

	slog.SetDefault(slog.New(newHandler).WithGroup("NewOuter").WithGroup("NewInner"))

	slog.Log(context.Background(), slog.Level(-8), "Meltdown")
	slog.Debug("Debugging", "Number", 100)
	slog.Warn("Warning", slog.Group("req",
		slog.String("url", "/api"),
		slog.Int("id", 42),
		"user", "name",
	))

	logger := slog.New(prettyslog.NewHandler(os.Stdout, nil))

	logger.Error("db init failed", "err", "invalid db url")

	logger.Info("requests",
		slog.Group("req",
			slog.String("url", "/api"),
			slog.Int("id", 42),
		),
		slog.Int("user", 12),
		slog.Group("reqs2",
			slog.String("url", "/api"),
			slog.Group("reqs4", "url", "/api", slog.Group("reqs5", "url", "/api", slog.Attr{
				Key:   "id",
				Value: slog.IntValue(70),
			})),
		),
		slog.Bool("bool", true),
	)
}
