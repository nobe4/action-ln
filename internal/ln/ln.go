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
	"github.com/nobe4/action-ln/internal/format"
	"github.com/nobe4/action-ln/internal/github"
	"github.com/nobe4/action-ln/internal/log"
)

func Run(ctx context.Context, e environment.Environment, g *github.GitHub) error {
	c, err := getConfig(ctx, g, e)
	if err != nil {
		return err
	}

	f := format.New(c, e)

	if err := c.Populate(ctx, g); err != nil {
		return fmt.Errorf("failed to populate config: %w", err)
	}

	groups := c.Links.Groups()

	if err := processGroups(ctx, g, f, groups); err != nil {
		return fmt.Errorf("failed to process the groups: %w", err)
	}

	return nil
}

func getConfig(ctx context.Context, g *github.GitHub, e environment.Environment) (*config.Config, error) {
	log.Group("Get config")

	log.Debug("Get config commit", "repo", e.Repo)

	b, err := g.GetDefaultBranch(ctx, e.Repo)
	if err != nil {
		return nil, fmt.Errorf("failed to get default branch: %w", err)
	}

	f := github.File{Repo: e.Repo, Path: e.Config, Commit: b.Commit.SHA, Ref: b.Name}

	log.Debug("Get config file", "file", f)

	if err := g.GetFile(ctx, &f); err != nil {
		return nil, fmt.Errorf("failed to get config %#v: %w", f, err)
	}

	log.Debug("Create config object", "default.repo", e.Repo)

	c := config.New()
	c.Defaults.Repo = e.Repo
	c.Source = f

	log.Debug("Parse config file", "sha", f.Commit)

	if err := c.Parse(strings.NewReader(f.Content)); err != nil {
		return nil, fmt.Errorf("failed to parse config %#v: %w", f, err)
	}

	log.Debug("Parsed config", "config", c)

	log.GroupEnd()

	return c, nil
}
