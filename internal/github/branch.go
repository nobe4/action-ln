package github

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
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

	if err := g.req(ctx, http.MethodGet, path, nil, &b); err != nil {
		return b, fmt.Errorf("failed to get branch: %w", err)
	}

	return b, nil
}

// https://docs.github.com/en/rest/git/refs?apiVersion=2022-11-28#create-a-reference
func (g GitHub) CreateBranch(ctx context.Context, repo Repo, branch, sha string) (Branch, error) {
	b := Branch{
		Name: branch,
		Commit: Commit{
			SHA: sha,
		},
	}

	path := fmt.Sprintf("/repos/%s/%s/git/refs", repo.Owner.Login, repo.Repo)

	body, err := json.Marshal(struct {
		Ref string `json:"ref"`
		SHA string `json:"sha"`
	}{
		Ref: "refs/heads/" + branch,
		SHA: sha,
	})
	if err != nil {
		return b, fmt.Errorf("failed to marshal request: %w", err)
	}

	if err := g.req(ctx, http.MethodPost, path, bytes.NewReader(body), nil); err != nil {
		return b, fmt.Errorf("failed to create branch: %w", err)
	}

	return b, nil
}
