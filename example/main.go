package main

import (
	"log/slog"
	"os"

	"github.com/kyr-go/prettyslog"
)

func main() {
	newHandler := prettyslog.NewHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelDebug,
	}).WithGroup("Outer").WithAttrs([]slog.Attr{slog.Group("test", "test1", 1, "test2", 2)})

	slog.SetDefault(slog.New(newHandler).With("hello", "world").WithGroup("Main").WithGroup("Sub"))

	slog.Debug("hello world", "test", "hey", "value", 36.44)

	slog.Warn("testing", slog.Group("req",
		slog.String("url", "/api"),
		slog.Int("id", 42),
		"user", "name",
	))

	logger := slog.New(prettyslog.NewHandler(os.Stdout, nil))

	logger.Error("db init failed", "test", "testing")

	logger.Info("requests",
		slog.Group("req",
			slog.String("url", "/api"),
			slog.Int("id", 42),
			"hey", 12,
		),
		slog.Int("user", 12),
		slog.Group("reqs2",
			slog.String("url", "/api"),
			slog.Int("id", 42),
			slog.Group("reqs3", "url", "/api", "heh", "hehe"),
			slog.Group("reqs4", "url", "/api", slog.Group("reqs5", "url", "/api", slog.Attr{
				Key:   "id",
				Value: slog.IntValue(70),
			})),
		),
		slog.Bool("hey", true),
	)
}
