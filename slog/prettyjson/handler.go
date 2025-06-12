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

func NewHandler(w io.Writer, slogOpts *slog.HandlerOptions, opts ...Option) (slog.Handler, error) {
	handlerOpts := &options{
		styleName: "monokai", // default style
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

	l := lexers.Get("json")
	if l == nil {
		return nil, fmt.Errorf("%w: failed to get lexer for json", ErrCreationFailed)
	}

	l = chroma.Coalesce(l)

	f := formatters.Get("terminal")
	if f == nil {
		return nil, fmt.Errorf("%w: failed to get formatter for terminal", ErrCreationFailed)
	}

	s := styles.Get(handlerOpts.styleName)
	if s == nil {
		return nil, fmt.Errorf("%w: failed to get style for %q", ErrCreationFailed, handlerOpts.styleName)
	}

	return &handler{
		h:           factory(w),
		m:           new(sync.Mutex),
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
	h slog.Handler
	m *sync.Mutex

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

	prettyBytes := pretty.Pretty(buf.Bytes())

	it, err := h.lexer.Tokenise(nil, string(prettyBytes))
	if err != nil {
		return err
	}

	h.m.Lock()
	defer h.m.Unlock()

	return h.formatter.Format(h.out, h.style, it)
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
