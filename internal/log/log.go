/*
https://pkg.go.dev/log/slog#pkg-constants
https://github.com/golang/example/blob/master/slog-handler-guide/README.md
https://github.com/actions/toolkit/blob/253e837c4db937cac18949bc65f0ffdd87496033/packages/core/src/command.ts
*/

package log

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"strings"
	"sync"
)

const (
	LevelDebug    = slog.LevelDebug
	LevelInfo     = slog.LevelInfo
	LevelNotice   = slog.Level(2)
	LevelWarn     = slog.LevelWarn
	LevelError    = slog.LevelError
	LevelGroup    = slog.Level(10)
	LevelGroupEnd = slog.Level(11)

	buflen = 1024
)

var errCannotWrite = errors.New("cannot write to output")

func Info(msg string, attrs ...any) {
	slog.Info(msg, attrs...)
}

func Debug(msg string, attrs ...any) {
	slog.Debug(msg, attrs...)
}

func Error(msg string, attrs ...any) {
	slog.Error(msg, attrs...)
}

func Warn(msg string, attrs ...any) {
	slog.Warn(msg, attrs...)
}

func Notice(msg string, attrs ...any) {
	slog.Log(context.Background(), LevelNotice, msg, attrs...)
}

func Group(name string) {
	slog.Log(context.Background(), LevelGroup, name)
}

func GroupEnd() {
	slog.Log(context.Background(), LevelGroupEnd, "")
}

type GitHubHandler struct {
	opts GitHubHandlerOptions
	mu   *sync.Mutex
	out  io.Writer
}

type GitHubHandlerOptions struct {
	Level slog.Leveler
}

func NewGitHubHandler(out io.Writer, debug bool) *GitHubHandler {
	h := &GitHubHandler{
		out: out,
		opts: GitHubHandlerOptions{
			Level: slog.LevelInfo,
		},
		mu: &sync.Mutex{},
	}

	if debug {
		h.opts.Level = slog.LevelDebug
	}

	return h
}

func (h *GitHubHandler) Enabled(_ context.Context, l slog.Level) bool {
	return l >= h.opts.Level.Level()
}

func (h *GitHubHandler) write(p []byte) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	_, err := h.out.Write(p)

	return fmt.Errorf("%w: %w", errCannotWrite, err)
}

// Format
// [::command key=value[,key=value]::]message.
func (h *GitHubHandler) Handle(_ context.Context, r slog.Record) error {
	command := ""

	switch r.Level {
	case LevelDebug:
		command = "debug"
	case LevelWarn:
		command = "warning"
	case LevelError:
		command = "error"
	case LevelNotice:
		command = "notice"

	// Info only prints the data, and doesn't add any context.
	case LevelInfo:
		return h.write([]byte(r.Message + "\n"))

	// Groups can be handled directly
	case LevelGroup:
		return h.write([]byte("::group::" + r.Message + "\n"))
	case LevelGroupEnd:
		return h.write([]byte("::groupend::\n"))
	}

	buf := make([]byte, 0, buflen)

	buf = fmt.Appendf(buf, "::%s", command)
	buf = fmt.Appendf(buf, "%s", h.formatAttrs(r))
	buf = fmt.Appendf(buf, "::%s\n", r.Message)

	return h.write(buf)
}

func (h *GitHubHandler) formatAttrs(r slog.Record) string {
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

func (h *GitHubHandler) WithAttrs(_ []slog.Attr) slog.Handler {
	// TODO: implement?
	return h
}

func (h *GitHubHandler) WithGroup(_ string) slog.Handler {
	// TODO: implement?
	return h
}
