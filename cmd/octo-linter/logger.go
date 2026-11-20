package main

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"slices"
	"strings"
)

type customHandler struct {
	opts  *slog.HandlerOptions
	attrs []slog.Attr
	group string
	w     io.Writer
}

func (h *customHandler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= h.opts.Level.Level()
}

func (h *customHandler) Handle(_ context.Context, r slog.Record) error {
	var buf strings.Builder

	fmt.Fprintf(&buf, "level=%s\n", r.Level.String())
	fmt.Fprintf(&buf, "msg=%q\n", r.Message)

	for _, a := range h.attrs {
		key := a.Key
		if h.group != "" {
			key = h.group + "." + a.Key
		}

		fmt.Fprintf(&buf, "%s=%q\n", key, a.Value.String())
	}

	r.Attrs(func(a slog.Attr) bool {
		fmt.Fprintf(&buf, "%s=%q\n", a.Key, a.Value.String())

		return true
	})
	buf.WriteString("---\n")

	// Errors from writing to the log writer are ignored because we cannot log a logging failure.
	_, _ = fmt.Fprint(h.w, buf.String())

	return nil
}

func (h *customHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	c := h.clone()
	c.attrs = append(c.attrs, attrs...)

	return c
}

func (h *customHandler) WithGroup(name string) slog.Handler {
	c := h.clone()
	c.group = name

	return c
}

func (h *customHandler) clone() *customHandler {
	return &customHandler{
		opts:  h.opts,
		attrs: slices.Clone(h.attrs),
		group: h.group,
		w:     h.w,
	}
}
