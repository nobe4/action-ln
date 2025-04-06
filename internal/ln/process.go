package ln

import (
	"context"
	"fmt"

	"github.com/nobe4/action-ln/internal/config"
	"github.com/nobe4/action-ln/internal/github"
	"github.com/nobe4/action-ln/internal/log"
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
	head, base, err := g.GetBaseAndHeadBranches(ctx, group[0].To.Repo, "test")
	if err != nil {
		return fmt.Errorf("failed to prepare branches: %w", err)
	}

	log.Debug("Parsed branches", "head", head, "base", base)

	for _, link := range group {
		log.Debug("Processing link", "link", link)

		needUpdate, err := link.NeedUpdate(ctx, g, head)
		if err != nil {
			return fmt.Errorf("failed to check if link needs update: %w", err)
		}

		log.Debug("Update needed", "", needUpdate)
	}

	return nil
}
