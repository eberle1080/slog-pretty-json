package prettyjson

import (
	"log/slog"
	"sync"
)

type Option func(*options)

// WithStyle sets the chroma style used for colorizing JSON.
// See https://xyproto.github.io/splash/docs/
func WithStyle(style string) Option {
	return func(opts *options) {
		opts.styleName = style
	}
}

// WithAttrs adds the provided attributes to every record handled.
func WithAttrs(attrs []slog.Attr) Option {
	return func(opts *options) {
		opts.attrs = attrs
	}
}

// WithGroup adds a slog group around all records.
func WithGroup(group string) Option {
	return func(opts *options) {
		opts.groups = append(opts.groups, group)
	}
}

// WithPretty toggles prettifying of the JSON output.
func WithPretty(pretty bool) Option {
	return func(opts *options) {
		opts.pretty = pretty
	}
}

// WithColor toggles ANSI color output.
func WithColor(color bool) Option {
	return func(opts *options) {
		opts.color = color
	}
}

type options struct {
	styleName string
	attrs     []slog.Attr
	groups    []string

	pretty bool
	color  bool

	mu *sync.Mutex
}

func (o *options) clone() *options {
	if o == nil {
		return nil
	}

	var attrs []slog.Attr

	if len(o.attrs) > 0 {
		attrs = make([]slog.Attr, len(o.attrs))
		copy(attrs, o.attrs)
	}

	var groups []string

	if len(o.groups) > 0 {
		groups = make([]string, len(o.groups))
		copy(groups, o.groups)
	}

	return &options{
		styleName: o.styleName,
		attrs:     attrs,
		groups:    groups,
		pretty:    o.pretty,
		color:     o.color,
		mu:        o.mu,
	}
}
