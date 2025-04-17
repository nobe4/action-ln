package ln

import (
	"context"
	"fmt"

	"github.com/nobe4/action-ln/internal/config"
	"github.com/nobe4/action-ln/internal/format"
	"github.com/nobe4/action-ln/internal/github"
	"github.com/nobe4/action-ln/internal/log"
)

func processGroups(ctx context.Context, g *github.GitHub, f format.Formatter, groups config.Groups) error {
	for id, l := range groups {
		log.Group("Processing group " + id)

		if err := processLinks(ctx, g, f, l); err != nil {
			return err
		}

		log.GroupEnd()
	}

	return nil
}

func processLinks(ctx context.Context, g *github.GitHub, f format.Formatter, l config.Links) error {
	toRepo := l[0].To.Repo

	base, head, err := g.GetBaseAndHeadBranches(ctx, toRepo, format.HeadBranch)
	if err != nil {
		return fmt.Errorf("failed to prepare branches: %w", err)
	}

	log.Debug("Parsed branches", "head", head, "base", base)

	updated, err := l.Update(ctx, g, head)
	if err != nil {
		return fmt.Errorf("failed to update the links: %w", err)
	}

	if !updated {
		// TODO don't create the PR, remove the branch if it's new.
		log.Debug("No link was updated, cleaning up...")
		log.Debug("head branch", "new", head.New)
	}

	pullBody, err := f.PullBody(l)
	if err != nil {
		return fmt.Errorf("failed to create pull request body: %w", err)
	}

	pull, err := g.GetOrCreatePull(ctx, toRepo, base.Name, head.Name, format.PullTitle, pullBody)
	if err != nil {
		return fmt.Errorf("failed to get pull request: %w", err)
	}

	log.Info("Created pull request", "pull", pull)

	return nil
}
