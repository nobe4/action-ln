package config

import (
	"context"
	"errors"
	"fmt"

	"github.com/nobe4/action-ln/internal/format"
	"github.com/nobe4/action-ln/internal/github"
	"github.com/nobe4/action-ln/internal/log"
)

const (
	commitMsgTemplate = `auto(ln): update {{ .Data.To.Path }}

Source: {{ .Environment.Server }}{{ .Data.From.HTMLPath }}
`
)

var (
	errMissingFrom = errors.New("from is missing")
	errMissingTo   = errors.New("to is missing")
	errInvalidFrom = errors.New("from is invalid")
	errInvalidTo   = errors.New("to is invalid")
)

type Link struct {
	From github.File `json:"from" yaml:"from"`
	To   github.File `json:"to"   yaml:"to"`
}

// The parsing can be done from a couple of various format, see ParseFile.
type RawLink struct {
	From any `yaml:"from"`
	To   any `yaml:"to"`
}

func (l *Link) String() string {
	return fmt.Sprintf("%s -> %s", l.From, l.To)
}

func (l *Link) Equal(other *Link) bool {
	return l.From.Equal(other.From) && l.To.Equal(other.To)
}

func (l *Link) NeedUpdate(ctx context.Context, g github.FileGetter, head github.Branch) (bool, error) {
	if l.From.Content == l.To.Content {
		log.Debug("Content is the same", "from", l.From, "to", l.To)

		return false, nil
	}

	headTo := &github.File{
		Repo: l.To.Repo,
		Path: l.To.Path,
		Ref:  head.Name,
	}

	log.Debug("Checking head content", "from", l.From, "to@head", headTo)

	if err := g.GetFile(ctx, headTo); err != nil {
		if errors.Is(err, github.ErrMissingFile) {
			log.Warn("File is missing", "to@head", headTo)

			return true, nil
		}

		return false, fmt.Errorf("failed to get to@head %s: %w", headTo, err)
	}

	if l.From.Content == headTo.Content {
		log.Debug("Content is the same", "from", l.From, "to@head", headTo)

		return false, nil
	}

	log.Debug("Content differs", "from", l.From, "to@head", headTo)

	return true, nil
}

func (l *Link) Update(ctx context.Context, g github.FileUpdater, f format.Formatter, head github.Branch) error {
	log.Info("Processing link", "link", l)

	l.To.Content = l.From.Content

	msg, err := f.Format(commitMsgTemplate, l)
	if err != nil {
		return fmt.Errorf("failed to format the commit message: %w", err)
	}

	newTo, err := g.UpdateFile(ctx, l.To, head.Name, msg)
	if err != nil {
		return fmt.Errorf("failed to update file: %w", err)
	}

	log.Info("Updated file", "new to", newTo)

	return nil
}

func (l *Link) populate(ctx context.Context, g github.FileGetter) error {
	if err := g.GetFile(ctx, &l.From); err != nil {
		return fmt.Errorf("%w %#v: %w", errMissingFrom, l.From, err)
	}

	if err := g.GetFile(ctx, &l.To); err != nil {
		if !errors.Is(err, github.ErrMissingFile) {
			return fmt.Errorf("%w %#v: %w", errMissingTo, l.To, err)
		}

		log.Debug("file does not exist", "file", l.To)
	}

	return nil
}

func (l *Link) fillMissing() {
	if l.To.Path == "" {
		l.To.Path = l.From.Path
	}

	if l.To.Repo.Empty() {
		l.To.Repo = l.From.Repo
	}
}
