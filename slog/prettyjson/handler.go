package prettyjson

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"sync"

	"github.com/alecthomas/chroma"
	"github.com/alecthomas/chroma/formatters"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
	"github.com/tidwall/pretty"
)

// NewHandler returns a slog.Handler that writes formatted JSON to w.
func NewHandler(w io.Writer, slogOpts *slog.HandlerOptions, opts ...Option) (slog.Handler, error) {
	handlerOpts := &options{
		styleName: "monokai", // default style
		pretty:    true,
		color:     true,
	}

	for _, opt := range opts {
		opt(handlerOpts)
	}

	return createHandler(w, slogOpts, handlerOpts)
}

func createHandler(w io.Writer, slogOpts *slog.HandlerOptions, handlerOpts *options) (slog.Handler, error) {
	factory := func(w io.Writer) slog.Handler {
		var h slog.Handler = slog.NewJSONHandler(w, slogOpts)

		if len(handlerOpts.attrs) > 0 {
			h = h.WithAttrs(handlerOpts.attrs)
		}

		if len(handlerOpts.groups) > 0 {
			for _, group := range handlerOpts.groups {
				h = h.WithGroup(group)
			}
		}

		return h
	}

	// If neither pretty nor color output is requested use the plain JSON handler.
	if !handlerOpts.pretty && !handlerOpts.color {
		return factory(w), nil
	}

	// Get the JSON lexer
	l := lexers.Get("json")
	if l == nil {
		return nil, fmt.Errorf("%w: failed to get lexer for json", ErrCreationFailed)
	}

	l = chroma.Coalesce(l)

	// Get the terminal formatter
	f := formatters.Get("terminal")
	if f == nil {
		return nil, fmt.Errorf("%w: failed to get formatter for terminal", ErrCreationFailed)
	}

	// Get the style
	s := styles.Get(handlerOpts.styleName)
	if s == nil {
		return nil, fmt.Errorf("%w: failed to get style for %q", ErrCreationFailed, handlerOpts.styleName)
	}

	return &handler{
		h:           factory(w),
		mu:          new(sync.Mutex),
		factory:     factory,
		out:         w,
		lexer:       l,
		formatter:   f,
		style:       s,
		slogOpts:    slogOpts,
		handlerOpts: handlerOpts,
	}, nil
}

type handler struct {
	h  slog.Handler
	mu *sync.Mutex

	factory func(w io.Writer) slog.Handler

	out io.Writer

	lexer     chroma.Lexer
	formatter chroma.Formatter
	style     *chroma.Style

	slogOpts    *slog.HandlerOptions
	handlerOpts *options
}

func (h *handler) Handle(ctx context.Context, record slog.Record) error {
	var buf bytes.Buffer
	handler := h.factory(&buf)

	if err := handler.Handle(ctx, record); err != nil {
		return err
	}

	var prettyBytes []byte
	if h.handlerOpts.pretty {
		prettyBytes = pretty.Pretty(buf.Bytes())
	} else {
		prettyBytes = buf.Bytes()
	}

	if !h.handlerOpts.color {
		_, err := h.out.Write(prettyBytes)
		return err
	}

	it, err := h.lexer.Tokenise(nil, string(prettyBytes))
	if err != nil {
		return err
	}

	var outBuf bytes.Buffer

	if err = h.formatter.Format(&outBuf, h.style, it); err != nil {
		return err
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	_, err = h.out.Write(outBuf.Bytes())

	return err
}

func (h *handler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.h.Enabled(ctx, level)
}

func (h *handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	newOpts := h.handlerOpts.clone()
	newOpts.attrs = append(newOpts.attrs, attrs...)

	handler, err := createHandler(h.out, h.slogOpts, newOpts)
	if err != nil {
		return slog.NewJSONHandler(h.out, h.slogOpts).WithAttrs(attrs) // fallback to JSON handler
	}

	return handler
}

func (h *handler) WithGroup(name string) slog.Handler {
	newOpts := h.handlerOpts.clone()
	newOpts.groups = append(newOpts.groups, name)

	handler, err := createHandler(h.out, h.slogOpts, newOpts)
	if err != nil {
		// Fallback to JSON handler
		var jh slog.Handler = slog.NewJSONHandler(h.out, h.slogOpts)

		for _, group := range h.handlerOpts.groups {
			jh = jh.WithGroup(group)
		}

		return jh
	}

	return handler
}
