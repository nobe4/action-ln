/*
Package ln is the main package for this codebase.

This is where the high-level logic is implemented.
*/
package ln

import (
	"context"
	"fmt"
	"strings"

	"github.com/nobe4/action-ln/internal/config"
	"github.com/nobe4/action-ln/internal/environment"
	"github.com/nobe4/action-ln/internal/github"
)

func Run(ctx context.Context, e environment.Environment, g *github.GitHub) error {
	c, err := getConfig(ctx, g, e)
	if err != nil {
		return err
	}

	if err := c.Populate(ctx, g); err != nil {
		return fmt.Errorf("failed to populate config: %w", err)
	}

	if err := processGroups(ctx, g, c.Links.Groups()); err != nil {
		return fmt.Errorf("failed to process the groups: %w", err)
	}

	return nil
}

func getConfig(ctx context.Context, g *github.GitHub, e environment.Environment) (*config.Config, error) {
	f := github.File{Repo: e.Repo, Path: e.Config}

	if err := g.GetFile(ctx, &f); err != nil {
		return nil, fmt.Errorf("failed to get config %#v: %w", f, err)
	}

	c := config.New()
	c.Defaults.Repo = e.Repo

	if err := c.Parse(strings.NewReader(f.Content)); err != nil {
		return nil, fmt.Errorf("failed to parse config %#v: %w", f, err)
	}

	return c, nil
}
