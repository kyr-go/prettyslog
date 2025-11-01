package prettyslog_test

import (
	"log/slog"
	"os"
	"testing"

	"github.com/kyr-go/prettyslog"
)

func BenchmarkLog(b *testing.B) {
	b.StopTimer()
	l := slog.New(prettyslog.NewHandler(os.Stdout, nil))

	for i := 0; i < b.N; i++ {
		b.StartTimer()
		l.Info("benchmarking", "i", i)
		b.StopTimer()
	}
}
