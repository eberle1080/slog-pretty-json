package prettyjson

import (
	"bytes"
	"context"
	"log/slog"
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
