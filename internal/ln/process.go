package ln

import (
	"context"
	"fmt"

	"github.com/nobe4/action-ln/internal/config"
	"github.com/nobe4/action-ln/internal/github"
	"github.com/nobe4/action-ln/internal/log"
)

const (
	branchName = "test"
)

func processGroups(ctx context.Context, g *github.GitHub, groups config.Groups) error {
	for id, group := range groups {
		log.Group("Processing group " + id)

		if err := processGroup(ctx, g, group); err != nil {
			return err
		}

		log.GroupEnd()
	}

	return nil
}

func processGroup(ctx context.Context, g *github.GitHub, group config.Links) error {
	base, head, err := g.GetBaseAndHeadBranches(ctx, group[0].To.Repo, branchName)
	if err != nil {
		return fmt.Errorf("failed to prepare branches: %w", err)
	}

	log.Debug("Parsed branches", "head", head, "base", base)

	for _, link := range group {
		if err := processLink(ctx, g, link, head); err != nil {
			return fmt.Errorf("failed to process link: %w", err)
		}
	}

	return nil
}

func processLink(ctx context.Context, g *github.GitHub, link *config.Link, head github.Branch) error {
	log.Debug("Processing link", "link", link)

	needUpdate, err := link.NeedUpdate(ctx, g, head)
	if err != nil {
		return fmt.Errorf("failed to check if link needs update: %w", err)
	}

	if !needUpdate {
		log.Debug("Update not needed")

		return nil
	}

	log.Debug("Update needed")

	link.To.Content = link.From.Content

	newTo, err := g.UpdateFile(ctx, link.To, head.Name, "test updating")
	if err != nil {
		return fmt.Errorf("failed to update file: %w", err)
	}

	log.Info("Updated file", "new to", newTo)

	return nil
}
