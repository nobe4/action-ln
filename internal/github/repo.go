package github

import (
	"context"
	"errors"
	"fmt"
	"net/http"
)

type Repo struct {
	Owner         User   `json:"owner"`
	Repo          string `json:"repo"`
	DefaultBranch string `json:"default_branch"`
}

var errGetRepo = errors.New("failed to get repo")

// https://docs.github.com/en/rest/repos/repos?apiVersion=2022-11-28
func (g *GitHub) GetDefaultBranch(ctx context.Context, repo Repo) (string, error) {
	path := fmt.Sprintf("/repos/%s/%s", repo.Owner.Login, repo.Repo)

	if _, err := g.req(ctx, http.MethodGet, path, nil, &repo); err != nil {
		return "", fmt.Errorf("%w: %w", errGetRepo, err)
	}

	return repo.DefaultBranch, nil
}
