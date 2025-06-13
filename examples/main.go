package main

import (
	"log/slog"
	"os"

	"github.com/eberle1080/slog-pretty-json/slog/prettyjson"
)

func main() {
	handler, err := prettyjson.NewHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelDebug,
	}, prettyjson.WithStyle("github"))
	if err != nil {
		panic(err)
	}

	logger := slog.New(handler)
	slog.SetDefault(logger)

	slog.Info("This is an info message",
		slog.String("key1", "value1"),
		slog.Int("key2", 42),
		slog.Bool("key3", true))
}
