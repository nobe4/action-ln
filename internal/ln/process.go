package ln

import (
	"context"
	"fmt"

	"github.com/nobe4/action-ln/internal/config"
	"github.com/nobe4/action-ln/internal/environment"
	"github.com/nobe4/action-ln/internal/github"
	"github.com/nobe4/action-ln/internal/log"
)

func processGroups(ctx context.Context, g *github.GitHub, e environment.Environment, c *config.Config) error {
	for id, l := range c.Links.Groups() {
		log.Group("Processing group " + id)

		if err := processLinks(ctx, g, c, e, l); err != nil {
			return err
		}

		log.GroupEnd()
	}

	return nil
}

func processLinks(
	ctx context.Context,
	g *github.GitHub,
	c *config.Config,
	e environment.Environment,
	l config.Links,
) error {
	toRepo := l[0].To.Repo

	base, head, err := g.GetBaseAndHeadBranches(ctx, toRepo, branchName)
	if err != nil {
		return fmt.Errorf("failed to prepare branches: %w", err)
	}

	log.Debug("Parsed branches", "head", head, "base", base)

	updated := false

	for _, link := range l {
		linkUpdated, err := processLink(ctx, g, link, head)
		if err != nil {
			return fmt.Errorf("failed to process link: %w", err)
		}

		updated = updated || linkUpdated
	}

	// TODO
	// if !updated {
	// 	log.Debug("No link was updated, cleaning up...")
	// }

	pullBody, err := pullRequestBody(l, c, e)
	if err != nil {
		return fmt.Errorf("failed to create pull request body: %w", err)
	}

	pull, err := g.GetOrCreatePull(ctx, toRepo, base.Name, head.Name, pullTitle, pullBody)
	if err != nil {
		return fmt.Errorf("failed to get pull request: %w", err)
	}

	log.Info("Created pull request", "pull", pull)

	return nil
}

func processLink(ctx context.Context, g *github.GitHub, link *config.Link, head github.Branch) (bool, error) {
	log.Info("Processing link", "link", link)

	needUpdate, err := link.NeedUpdate(ctx, g, head)
	if err != nil {
		return false, fmt.Errorf("failed to check if link needs update: %w", err)
	}

	if !needUpdate {
		log.Debug("Update not needed")

		return false, nil
	}

	log.Debug("Update needed")

	link.To.Content = link.From.Content

	newTo, err := g.UpdateFile(ctx, link.To, head.Name, "test updating")
	if err != nil {
		return false, fmt.Errorf("failed to update file: %w", err)
	}

	log.Info("Updated file", "new to", newTo)

	return true, nil
}
