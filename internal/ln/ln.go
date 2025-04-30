/*
Package ln is the main package for this codebase.

This is where the high-level logic is implemented.
*/
package ln

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/nobe4/action-ln/internal/config"
	"github.com/nobe4/action-ln/internal/environment"
	contextfmt "github.com/nobe4/action-ln/internal/format/context"
	"github.com/nobe4/action-ln/internal/github"
	"github.com/nobe4/action-ln/internal/log"
)

func Run(ctx context.Context, e environment.Environment, g *github.GitHub) error {
	c, err := getConfig(ctx, g, e)
	if err != nil {
		return err
	}

	f := contextfmt.New(c, e)

	groups := c.Links.Groups()

	log.Debug("Processing groups", "groups", "\n"+groups.String())

	if err := processGroups(ctx, g, f, groups); err != nil {
		return fmt.Errorf("failed to process the groups: %w", err)
	}

	return nil
}

func getConfig(ctx context.Context, g *github.GitHub, e environment.Environment) (*config.Config, error) {
	f, err := readConfig(ctx, g, e)
	if err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	c := config.New()
	c.Source = f

	if err := c.Parse(strings.NewReader(f.Content)); err != nil {
		return nil, fmt.Errorf("failed to parse config %#v: %w", f, err)
	}

	if err := c.Populate(ctx, g); err != nil {
		return nil, fmt.Errorf("failed to populate config: %w", err)
	}

	log.Debug("Parsed config", "config", c)

	return c, nil
}

func readConfig(ctx context.Context, g *github.GitHub, e environment.Environment) (github.File, error) {
	log.Group("Read config")
	defer log.GroupEnd()

	if e.LocalConfig == "" {
		return readConfigFromGitHub(ctx, g, e)
	}

	return readConfigFromFS(e.LocalConfig)
}

func readConfigFromGitHub(ctx context.Context, g *github.GitHub, e environment.Environment) (github.File, error) {
	log.Info("Read config from GitHub", "repo", e.Repo)

	b, err := g.GetDefaultBranch(ctx, e.Repo)
	if err != nil {
		return github.File{}, fmt.Errorf("failed to get default branch: %w", err)
	}

	f := github.File{Repo: e.Repo, Path: e.Config, Commit: b.Commit.SHA, Ref: b.Name}

	log.Info("Get config file", "file", f)

	if err := g.GetFile(ctx, &f); err != nil {
		return github.File{}, fmt.Errorf("failed to get config %#v: %w", f, err)
	}

	return f, nil
}

func readConfigFromFS(path string) (github.File, error) {
	log.Info("Read config from Filesystem", "path", path)

	content, err := os.ReadFile(path)
	if err != nil {
		return github.File{}, fmt.Errorf("failed to read config file %s: %w", path, err)
	}

	return github.File{
		Content: string(content),
		Name:    filepath.Base(path),
		Path:    path,
		HTMLURL: "file://" + path,
		Commit:  "local_commit",
		Ref:     "local_ref",
		SHA:     "local_sha",

		Repo: github.Repo{
			Owner: github.User{Login: "local_owner"},
			Repo:  "local_repo",
		},
	}, nil
}
