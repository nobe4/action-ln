package github

import (
	"context"
	"fmt"
)

type Commit struct {
	SHA string `json:"sha"`
}

type Branch struct {
	Name   string `json:"name"`
	Commit Commit `json:"commit"`
}

// https://docs.github.com/en/rest/branches/branches?apiVersion=2022-11-28#get-a-branch
func (g GitHub) GetBranch(ctx context.Context, repo Repo, branch string) (Branch, error) {
	b := Branch{}

	path := fmt.Sprintf("/repos/%s/%s/branches/%s", repo.Owner.Login, repo.Repo, branch)

	if err := g.req(ctx, "GET", path, nil, &b); err != nil {
		return b, fmt.Errorf("failed to get branch: %w", err)
	}

	return b, nil
}
