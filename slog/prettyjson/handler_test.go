package prettyjson

import (
	"bytes"
	"context"
	"log/slog"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestHandlerWritesOutput(t *testing.T) {
	var buf bytes.Buffer
	h, err := NewHandler(&buf, nil)
	if err != nil {
		t.Fatalf("new handler: %v", err)
	}

	r := slog.NewRecord(time.Now(), slog.LevelInfo, "msg", 0)
	r.AddAttrs(slog.String("k", "v"))

	if err := h.Handle(context.Background(), r); err != nil {
		t.Fatalf("handle: %v", err)
	}
	if buf.Len() == 0 {
		t.Fatal("no output written")
	}
}

func TestHandlerEnabled(t *testing.T) {
	var buf bytes.Buffer
	h, err := NewHandler(&buf, &slog.HandlerOptions{Level: slog.LevelWarn})
	if err != nil {
		t.Fatalf("new handler: %v", err)
	}

	if !h.Enabled(context.Background(), slog.LevelError) {
		t.Error("handler should be enabled for error level")
	}

	if h.Enabled(context.Background(), slog.LevelInfo) {
		t.Error("handler should not be enabled for info level")
	}
}

func TestHandlerWithAttrs(t *testing.T) {
	var buf bytes.Buffer
	h, err := NewHandler(&buf, nil)
	if err != nil {
		t.Fatalf("new handler: %v", err)
	}

	attrs := []slog.Attr{slog.String("test", "value")}
	h2 := h.WithAttrs(attrs)

	r := slog.NewRecord(time.Now(), slog.LevelInfo, "msg", 0)
	if err := h2.Handle(context.Background(), r); err != nil {
		t.Fatalf("handle: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "test") || !strings.Contains(output, "value") {
		t.Error("output should contain added attributes")
	}
}

func TestHandlerWithGroup(t *testing.T) {
	var buf bytes.Buffer
	h, err := NewHandler(&buf, nil)
	if err != nil {
		t.Fatalf("new handler: %v", err)
	}

	h2 := h.WithGroup("testgroup")

	r := slog.NewRecord(time.Now(), slog.LevelInfo, "msg", 0)
	r.AddAttrs(slog.String("k", "v"))
	if err := h2.Handle(context.Background(), r); err != nil {
		t.Fatalf("handle: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "testgroup") {
		t.Error("output should contain group name")
	}
}

func TestHandlerConcurrency(t *testing.T) {
	var buf bytes.Buffer
	h, err := NewHandler(&buf, nil)
	if err != nil {
		t.Fatalf("new handler: %v", err)
	}

	var wg sync.WaitGroup
	numGoroutines := 100

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			r := slog.NewRecord(time.Now(), slog.LevelInfo, "msg", 0)
			r.AddAttrs(slog.Int("goroutine", i))
			if err := h.Handle(context.Background(), r); err != nil {
				t.Errorf("handle failed: %v", err)
			}
		}(i)
	}

	wg.Wait()
}

func TestWithStyle(t *testing.T) {
	var buf bytes.Buffer
	h, err := NewHandler(&buf, nil, WithStyle("github"))
	if err != nil {
		t.Fatalf("new handler: %v", err)
	}

	r := slog.NewRecord(time.Now(), slog.LevelInfo, "msg", 0)
	if err := h.Handle(context.Background(), r); err != nil {
		t.Fatalf("handle: %v", err)
	}
}

func TestWithStyleEmpty(t *testing.T) {
	var buf bytes.Buffer
	h, err := NewHandler(&buf, nil, WithStyle(""))
	if err != nil {
		t.Fatalf("new handler with empty style should not fail: %v", err)
	}

	r := slog.NewRecord(time.Now(), slog.LevelInfo, "msg", 0)
	if err := h.Handle(context.Background(), r); err != nil {
		t.Fatalf("handle: %v", err)
	}
}

func TestWithPretty(t *testing.T) {
	var buf bytes.Buffer
	h, err := NewHandler(&buf, nil, WithPretty(false))
	if err != nil {
		t.Fatalf("new handler: %v", err)
	}

	r := slog.NewRecord(time.Now(), slog.LevelInfo, "msg", 0)
	if err := h.Handle(context.Background(), r); err != nil {
		t.Fatalf("handle: %v", err)
	}
}

func TestWithColor(t *testing.T) {
	var buf bytes.Buffer
	h, err := NewHandler(&buf, nil, WithColor(false))
	if err != nil {
		t.Fatalf("new handler: %v", err)
	}

	r := slog.NewRecord(time.Now(), slog.LevelInfo, "msg", 0)
	if err := h.Handle(context.Background(), r); err != nil {
		t.Fatalf("handle: %v", err)
	}
}

func TestWithAttrsOption(t *testing.T) {
	var buf bytes.Buffer
	attrs := []slog.Attr{slog.String("global", "value")}
	h, err := NewHandler(&buf, nil, WithAttrs(attrs))
	if err != nil {
		t.Fatalf("new handler: %v", err)
	}

	r := slog.NewRecord(time.Now(), slog.LevelInfo, "msg", 0)
	if err := h.Handle(context.Background(), r); err != nil {
		t.Fatalf("handle: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "global") || !strings.Contains(output, "value") {
		t.Error("output should contain global attributes")
	}
}

func TestWithGroupOption(t *testing.T) {
	var buf bytes.Buffer
	h, err := NewHandler(&buf, nil, WithGroup("global"))
	if err != nil {
		t.Fatalf("new handler: %v", err)
	}

	r := slog.NewRecord(time.Now(), slog.LevelInfo, "msg", 0)
	r.AddAttrs(slog.String("k", "v"))
	if err := h.Handle(context.Background(), r); err != nil {
		t.Fatalf("handle: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "global") {
		t.Error("output should contain global group")
	}
}
