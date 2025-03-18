/*
Package log provides a simple logging interface.

It sets up the levels and adds handy wrappers that will work with any handlers.

Refs:
- https://docs.github.com/en/actions/writing-workflows/choosing-what-your-workflow-does/workflow-commands-for-github-actions
- https://pkg.go.dev/log/slog#pkg-constants
- https://github.com/golang/example/blob/master/slog-handler-guide/README.md
- https://github.com/actions/toolkit/blob/253e837c4db937cac18949bc65f0ffdd87496033/packages/core/src/command.ts
*/
package log

import (
	"context"
	"errors"
	"log/slog"
)

var ErrCannotWrite = errors.New("cannot write to output")

const (
	LevelDebug    = slog.LevelDebug
	LevelInfo     = slog.LevelInfo
	LevelNotice   = slog.Level(2)
	LevelWarn     = slog.LevelWarn
	LevelError    = slog.LevelError
	LevelGroup    = slog.Level(10)
	LevelGroupEnd = slog.Level(11)
)

type Options struct {
	Level slog.Leveler
}

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
