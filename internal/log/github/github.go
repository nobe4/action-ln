/*
Package github provides a log handler that formats according to Github Actions
specifications.

Refs:
- https://docs.github.com/en/actions/writing-workflows/choosing-what-your-workflow-does/workflow-commands-for-github-actions
- https://pkg.go.dev/log/slog#pkg-constants
- https://github.com/golang/example/blob/master/slog-handler-guide/README.md
- https://github.com/actions/toolkit/blob/253e837c4db937cac18949bc65f0ffdd87496033/packages/core/src/command.ts
*/
package github

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"strings"
	"sync"

	"github.com/nobe4/action-ln/internal/log"
)

const buflen = 1024

type Handler struct {
	opts log.Options
	mu   *sync.Mutex
	out  io.Writer
}

func New(out io.Writer, o log.Options) *Handler {
	h := &Handler{
		out:  out,
		opts: o,
		mu:   &sync.Mutex{},
	}

	return h
}

func (h *Handler) Enabled(_ context.Context, l slog.Level) bool {
	return l >= h.opts.Level.Level()
}

func (h *Handler) write(p []byte) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	_, err := h.out.Write(p)

	return fmt.Errorf("%w: %w", log.ErrCannotWrite, err)
}

func (h *Handler) Handle(_ context.Context, r slog.Record) error {
	command := ""

	switch r.Level {
	// Info can be handled directly
	case log.LevelInfo:
		return h.write([]byte(r.Message + h.formatAttrs(r) + "\n"))

	// Groups can be handled directly
	case log.LevelGroup:
		return h.write([]byte("::group::" + r.Message + "\n"))
	case log.LevelGroupEnd:
		return h.write([]byte("::groupend::\n"))

	case log.LevelDebug:
		command = "debug"
	case log.LevelWarn:
		command = "warning"
	case log.LevelError:
		command = "error"
	case log.LevelNotice:
		command = "notice"
	}

	buf := make([]byte, 0, buflen)

	buf = fmt.Appendf(buf, "::%s", command)
	buf = fmt.Appendf(buf, "%s", h.formatAttrs(r))
	buf = fmt.Appendf(buf, "::%s\n", r.Message)

	return h.write(buf)
}

func (*Handler) formatAttrs(r slog.Record) string {
	attrs := []string{}

	r.Attrs(func(a slog.Attr) bool {
		attrs = append(attrs, fmt.Sprintf("%s=%s", a.Key, a.Value))

		return true
	})

	if len(attrs) > 0 {
		return " " + strings.Join(attrs, ",")
	}

	return ""
}

func (h *Handler) WithAttrs(_ []slog.Attr) slog.Handler {
	// TODO: implement?
	return h
}

func (h *Handler) WithGroup(_ string) slog.Handler {
	// TODO: implement?
	return h
}
