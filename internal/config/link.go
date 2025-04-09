package config

import (
	"context"
	"errors"
	"fmt"

	"github.com/nobe4/action-ln/internal/github"
	"github.com/nobe4/action-ln/internal/log"
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

func (l *Link) String() string {
	return fmt.Sprintf("%s -> %s", l.From, l.To)
}

func (l *Link) Equal(other *Link) bool {
	return l.From.Equal(other.From) && l.To.Equal(other.To)
}

func (l *Link) NeedUpdate(ctx context.Context, g github.FileGetter, head github.Branch) (bool, error) {
	// TODO: if the content is equal, this is not needed.
	if head.New {
		log.Debug("Head is new", "head", head)

		return true, nil
	}

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

type Links []*Link

type Groups map[string]Links

func (l Links) Groups() Groups {
	g := make(Groups)

	for _, link := range l {
		g[link.To.Repo.String()] = append(g[link.To.Repo.String()], link)
	}

	return g
}

type RawLink struct {
	From any `yaml:"from"`
	To   any `yaml:"to"`
}

func (c *Config) parseLinks(raw []RawLink) (Links, error) {
	links := Links{}

	for i, rl := range raw {
		log.Debug("Parse link", "index", i, "raw", rl)

		l, err := c.parseLink(rl)
		if err != nil {
			return nil, err
		}

		links = append(links, l)
	}

	return links, nil
}

func (c *Config) parseLink(raw RawLink) (*Link, error) {
	from, err := c.parseFile(raw.From)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", errInvalidFrom, err)
	}

	to, err := c.parseFile(raw.To)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", errInvalidTo, err)
	}

	return &Link{From: from, To: to}, nil
}
