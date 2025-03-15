package ln

import (
	"context"
	"fmt"
	"os"

	"github.com/nobe4/action-ln/internal/config"
	"github.com/nobe4/action-ln/internal/github"
)

func processGroups(ctx context.Context, g *github.GitHub, groups config.Groups) error {
	for id, group := range groups {
		fmt.Fprintf(os.Stderr, "\ngroup: %s\n", id)

		if err := processGroup(ctx, g, group); err != nil {
			return err
		}
	}

	return nil
}

func processGroup(ctx context.Context, g *github.GitHub, group config.Links) error {
	head, base, err := g.GetBaseAndHeadBranches(ctx, group[0].To.Repo, "test")
	if err != nil {
		return fmt.Errorf("failed to prepare branches: %w", err)
	}

	fmt.Fprintf(os.Stderr, "  base branch: %v\n", base)
	fmt.Fprintf(os.Stderr, "  head branch: %v\n", head)

	for _, link := range group {
		fmt.Fprintf(os.Stderr, "  link: %s\n", link)

		needUpdate, err := link.NeedUpdate(ctx, g, head)
		if err != nil {
			return fmt.Errorf("failed to check if link needs update: %w", err)
		}

		fmt.Fprintf(os.Stderr, "    Need Update: %v\n", needUpdate)
	}

	return nil
}
