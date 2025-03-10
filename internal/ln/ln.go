/*
Package ln is the main package for this codebase.

This is where the high-level logic is implemented.
*/
package ln

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/nobe4/action-ln/internal/config"
	"github.com/nobe4/action-ln/internal/environment"
	"github.com/nobe4/action-ln/internal/github"
)

func Run(ctx context.Context, e environment.Environment, g github.GitHub) error {
	c, err := getConfig(ctx, g, e)
	if err != nil {
		return err
	}

	fmt.Fprintf(os.Stdout, "Configuration before: %s\n", c)

	if err := populateConfig(ctx, &c, g, e); err != nil {
		return err
	}

	fmt.Fprintf(os.Stdout, "Configuration after: %s\n", c)

	return nil
}

func getConfig(ctx context.Context, g github.GitHub, e environment.Environment) (config.Config, error) {
	f := github.File{Repo: e.Repo, Path: e.Config}

	if err := g.GetFile(ctx, &f); err != nil {
		return config.Config{}, fmt.Errorf("failed to get config %#v: %w", f, err)
	}

	c, err := config.Parse(strings.NewReader(f.Content))
	if err != nil {
		return config.Config{}, fmt.Errorf("failed to parse config %#v: %w", f, err)
	}

	return c, nil
}

func populateConfig(ctx context.Context, c *config.Config, g github.GitHub, e environment.Environment) error {
	for i, l := range c.Links {
		if err := populateLink(ctx, &l, g, e); err != nil {
			return fmt.Errorf("failed to populate link %#v: %w", l, err)
		}

		c.Links[i] = l
	}

	return nil
}

func populateLink(ctx context.Context, l *config.Link, g github.GitHub, e environment.Environment) error {
	if l.From.Repo.Empty() {
		l.From.Repo = e.Repo
	}

	if err := g.GetFile(ctx, &l.From); err != nil {
		return fmt.Errorf("failed to get from file %#v: %w", l.From, err)
	}

	if l.To.Repo.Empty() {
		l.To.Repo = e.Repo
	}

	if err := g.GetFile(ctx, &l.To); err != nil {
		if !errors.Is(err, github.ErrMissingFile) {
			return fmt.Errorf("failed to get to file %#v: %w", l.To, err)
		}

		fmt.Fprintf(os.Stdout, "File %#v does not exist\n", l.To)
	}

	return nil
}
