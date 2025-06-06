package config

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/nobe4/action-ln/internal/format"
	"github.com/nobe4/action-ln/internal/github"
	"github.com/nobe4/action-ln/internal/log"
)

const (
	commitMsgTemplate = `auto(ln): update {{ .Data.To.Path }}

Source: {{ .Data.From.HTMLURL }}
`
	linkStringPartCount = 2
)

var (
	errMissingFrom       = errors.New("from is missing")
	errGettingRepo       = errors.New("failed to get repo")
	errMissingTo         = errors.New("to is missing")
	errInvalidFrom       = errors.New("from is invalid")
	errInvalidTo         = errors.New("to is invalid")
	errInvalidLinkFormat = errors.New("link format invalid, want 'from -> to'")
)

type Link struct {
	From github.File `json:"from" yaml:"from"`
	To   github.File `json:"to"   yaml:"to"`

	Status Status `json:"status" yaml:"status"`
}

type Status string

const (
	StatusFailedToCheck   Status = "failed to check for update"
	StatusFailedToUpdate  Status = "failed to update"
	StatusUpdateNotNeeded Status = "update not needed"
	StatusUpdated         Status = "updated"
)

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

func (l *Link) NeedUpdate(ctx context.Context, g github.Getter, head github.Branch) (bool, error) {
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

func (l *Link) Update(ctx context.Context, g github.Updater, f format.Formatter, head github.Branch) error {
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

func (c *Config) ParseLinkString(s string) (Link, error) {
	if s == "" {
		return Link{}, nil
	}

	p := strings.Split(s, " -> ")

	if l := len(p); l != linkStringPartCount {
		return Link{}, fmt.Errorf("%w, got %d for %q", errInvalidLinkFormat, l, s)
	}

	from, err := c.parseString(p[0])
	if err != nil {
		return Link{}, fmt.Errorf("%w %q: %w", errInvalidFrom, p[0], err)
	}

	to, err := c.parseString(p[1])
	if err != nil {
		return Link{}, fmt.Errorf("%w %q: %w", errInvalidTo, p[0], err)
	}

	return Link{From: from[0], To: to[0]}, nil
}

func (l *Link) populate(ctx context.Context, g github.Getter) error {
	if err := l.populateFrom(ctx, g); err != nil {
		return err
	}

	return l.populateTo(ctx, g)
}

func (l *Link) populateFrom(ctx context.Context, g github.Getter) error {
	// NOTE: Technically speaking, having the `Ref` is not needed to get the
	// content on the default branch. However, there's no way to get it from
	// `GetFile`, so getting it in advance is nicer for displaying it later.
	if l.From.Ref == "" {
		if err := g.GetRepo(ctx, &l.From.Repo); err != nil {
			return fmt.Errorf("%w %#v: %w", errGettingRepo, l.From, err)
		}

		l.From.Ref = l.From.Repo.DefaultBranch
	}

	if err := g.GetFile(ctx, &l.From); err != nil {
		return fmt.Errorf("%w %#v: %w", errMissingFrom, l.From, err)
	}

	return nil
}

func (l *Link) populateTo(ctx context.Context, g github.Getter) error {
	refs := []string{"auto-action-ln", l.To.Ref}

	for _, ref := range refs {
		l.To.Ref = ref

		err := g.GetFile(ctx, &l.To)
		if err == nil {
			return nil
		}

		if errors.Is(err, github.ErrMissingFile) {
			log.Debug("file does not exist", "file", l.To, "ref", l.To.Ref)

			continue
		}

		return fmt.Errorf("%w %#v: %w", errMissingTo, l.To, err)
	}

	return nil
}
