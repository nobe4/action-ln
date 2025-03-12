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
	headBranch, baseBranch, err := prepareBranches(ctx, g, group[0].To.Repo)
	if err != nil {
		return fmt.Errorf("failed to prepare branches: %w", err)
	}

	fmt.Fprintf(os.Stderr, "  base branch: %v\n", baseBranch)
	fmt.Fprintf(os.Stderr, "  head branch: %v\n", headBranch)

	for _, link := range group {
		fmt.Fprintf(os.Stderr, "  link: %s\n", link)
		fmt.Fprintf(os.Stderr, "    Need Update: %v\n", link.NeedsUpdate())
	}

	return nil
}

func prepareBranches(
	ctx context.Context,
	g *github.GitHub,
	r github.Repo,
) (base github.Branch, head github.Branch, err error) {
	// TODO: see if this can come with the Repo
	defaultBranchName, err := g.GetDefaultBranch(ctx, r)
	if err != nil {
		return base, head, fmt.Errorf("failed to get default branch name: %w", err)
	}

	base, err = g.GetBranch(ctx, r, defaultBranchName)
	if err != nil {
		return base, head, fmt.Errorf("failed to get default branch: %w", err)
	}

	head, err = g.GetOrCreateBranch(ctx, r, "test", base.Commit.SHA)
	if err != nil {
		return base, head, fmt.Errorf("failed to get default branch: %w", err)
	}

	return base, head, nil
}
