package prettyjson

import "log/slog"

type Option func(*options)

// WithStyle sets the style.
// See https://xyproto.github.io/splash/docs/
func WithStyle(style string) Option {
	return func(opts *options) {
		opts.styleName = style
	}
}

func WithAttrs(attrs []slog.Attr) Option {
	return func(opts *options) {
		opts.attrs = attrs
	}
}

func WithGroup(group string) Option {
	return func(opts *options) {
		opts.groups = append(opts.groups, group)
	}
}

type options struct {
	styleName string
	attrs     []slog.Attr
	groups    []string
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
	}
}
