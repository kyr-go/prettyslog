# Pretty Slog - A Prettier log/slog
[![Go Reference](https://pkg.go.dev/badge/github.com/kyr-go/prettyslog.svg)](https://pkg.go.dev/github.com/kyr-go/prettyslog)
[![Go](https://github.com/kyr-go/prettyslog/actions/workflows/go.yml/badge.svg)](https://github.com/kyr-go/prettyslog/actions/workflows/go.yml)

<img width="430" height="93" alt="image" src="https://github.com/user-attachments/assets/79a708a2-0165-435c-a51e-84d51130f417" />

## Install
```bash
go get github.com/kyr-go/prettyslog@latest
```

## Basic Usage
```go
package main

import (
	"log/slog"
	"os"

	"github.com/kyr-go/prettyslog"
)

func main() {
	opts := slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelDebug,
	}
	slog.SetDefault(slog.New(prettyslog.NewHandler(os.Stdout, &opts)))

	slog.Debug("Debug Message")
	slog.Info("Hello World")
	slog.Warn("Warning Message")
	slog.Error("Error Message")
}
```
