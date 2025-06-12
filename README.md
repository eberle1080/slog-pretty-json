# slog-pretty-json

A simple [slog](https://pkg.go.dev/log/slog) handler that outputs nicely
formatted JSON. The handler can be used as a drop in replacement for the
standard `slog.JSONHandler` while still supporting colorised output.

Uses [github.com/tidwall/pretty](https://github.com/tidwall/pretty) for formatting and
[github.com/alecthomas/chroma](https://github.com/alecthomas/chroma) for syntax highlighting.

![Example output](https://github.com/user-attachments/assets/89e91a50-15be-43f1-9973-0f4aaebc9f75)

## Installation

Add the module to your project:

```bash
go get github.com/eberle1080/slog-pretty-json/slog/prettyjson
```

## Basic usage

Create a handler and set it as the default logger. The example below mirrors
`testing/main.go` in this repository.

```go
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
```
