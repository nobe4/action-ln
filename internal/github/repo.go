package github

import (
	"context"
	"fmt"
	"net/http"
)

type Repo struct {
	Owner         User   `json:"owner"`
	Repo          string `json:"repo"`
	DefaultBranch string `json:"default_branch"`
}

// https://docs.github.com/en/rest/repos/repos?apiVersion=2022-11-28
func (g GitHub) GetDefaultBranch(ctx context.Context, repo Repo) (string, error) {
	path := fmt.Sprintf("/repos/%s/%s", repo.Owner.Login, repo.Repo)

	if _, err := g.req(ctx, http.MethodGet, path, nil, &repo); err != nil {
		// TODO: make constant error
		return "", fmt.Errorf("failed to get repo: %w", err)
	}

	return repo.DefaultBranch, nil
}
