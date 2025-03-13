package config

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/nobe4/action-ln/internal/github"
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
	if head.New {
		return true, nil
	}

	if l.From.Content == l.To.Content {
		return false, nil
	}

	headTo := &github.File{
		Repo: l.To.Repo,
		Path: l.To.Path,
		Ref:  head.Name,
	}

	if err := g.GetFile(ctx, headTo); err != nil {
		if errors.Is(err, github.ErrMissingFile) {
			fmt.Fprintln(os.Stdout, "    headTo: MISSING")

			return true, nil
		}

		return false, fmt.Errorf("failed to get file on branch: %w", err)
	}

	fmt.Fprintf(os.Stderr, "    headTo: %s\n", headTo)

	return l.From.Content != headTo.Content, nil
}

func (l *Link) populate(ctx context.Context, g github.FileGetter) error {
	if err := g.GetFile(ctx, &l.From); err != nil {
		return fmt.Errorf("%w %#v: %w", errMissingFrom, l.From, err)
	}

	if err := g.GetFile(ctx, &l.To); err != nil {
		if !errors.Is(err, github.ErrMissingFile) {
			return fmt.Errorf("%w %#v: %w", errMissingTo, l.To, err)
		}

		// TODO: make this a debug statement
		fmt.Fprintf(os.Stdout, "File %#v does not exist\n", l.To)
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

	for _, rl := range raw {
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
