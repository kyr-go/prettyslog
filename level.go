package prettyslog

import (
	"log/slog"
)

var levels = map[slog.Level]string{
	slog.LevelDebug: "\x1b[1;97;44mDEBUG\x1b[0m", // Blue background
	slog.LevelInfo:  "\x1b[1;97;42mINFO \x1b[0m", // Green background
	slog.LevelWarn:  "\x1b[1;97;43mWARN \x1b[0m", // Yellow background
	slog.LevelError: "\x1b[1;97;41mERROR\x1b[0m", // Red background
}
