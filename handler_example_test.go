package prettyslog_test

import (
	"log/slog"
	"os"

	"github.com/kyr-go/prettyslog"
)

func Example() {
	opts := slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelDebug,
	}
	slog.SetDefault(slog.New(prettyslog.NewHandler(os.Stdout, &opts)))

	slog.Debug("Debug Message")
	slog.Info("Info Message")
	slog.Warn("Warning Message")
	slog.Error("Error Message")
}
