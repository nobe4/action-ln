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
)

type RawLink struct {
	From any `yaml:"from"`
	To   any `yaml:"to"`
}

type Link struct {
	From github.File `json:"from" yaml:"from"`
	To   github.File `json:"to"   yaml:"to"`
}

func (c *Config) parseLinks(raw []RawLink) ([]Link, error) {
	links := []Link{}

	for _, rl := range raw {
		l, err := c.parseLink(rl)
		if err != nil {
			return nil, err
		}

		links = append(links, l)
	}

	return links, nil
}

func (c *Config) parseLink(raw RawLink) (Link, error) {
	from, err := parseFile(raw.From)
	if err != nil {
		return Link{}, err
	}

	to, err := parseFile(raw.To)
	if err != nil {
		return Link{}, err
	}

	// TODO: Move this into the file parsing
	if from.Repo.Empty() {
		from.Repo = c.Defaults.Repo
	}

	if to.Repo.Empty() {
		to.Repo = c.Defaults.Repo
	}

	return Link{From: from, To: to}, nil
}

func (l *Link) Populate(ctx context.Context, g github.FileGetter) error {
	if err := g.GetFile(ctx, &l.From); err != nil {
		return fmt.Errorf("%w %#v: %w", errMissingFrom, l.From, err)
	}

	if err := g.GetFile(ctx, &l.To); err != nil {
		if !errors.Is(err, github.ErrMissingFile) {
			return fmt.Errorf("%w %#v: %w", errMissingTo, l.To, err)
		}

		fmt.Fprintf(os.Stdout, "File %#v does not exist\n", l.To)
	}

	return nil
}
