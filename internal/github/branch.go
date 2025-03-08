package github

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

var (
	ErrNoBranch     = errors.New("branch not found")
	ErrGetBranch    = errors.New("failed to get branch")
	ErrCreateBranch = errors.New("failed to create branch")
	ErrBranchExists = errors.New("branch already exist")
)

type Commit struct {
	SHA string `json:"sha"`
}

type Branch struct {
	Name   string `json:"name"`
	Commit Commit `json:"commit"`
	New    bool   `json:"new"`
}

// https://docs.github.com/en/rest/branches/branches?apiVersion=2022-11-28#get-a-branch
func (g GitHub) GetBranch(ctx context.Context, repo Repo, branch string) (Branch, error) {
	b := Branch{}

	path := fmt.Sprintf("/repos/%s/%s/branches/%s", repo.Owner.Login, repo.Repo, branch)

	if status, err := g.req(ctx, http.MethodGet, path, nil, &b); err != nil {
		if status == http.StatusNotFound {
			return b, ErrNoBranch
		}

		return b, fmt.Errorf("%w: %w", ErrGetBranch, err)
	}

	b.New = false

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
		return b, fmt.Errorf("%w: %w", ErrMarshalRequest, err)
	}

	if status, err := g.req(ctx, http.MethodPost, path, bytes.NewReader(body), nil); err != nil {
		if status == http.StatusUnprocessableEntity {
			return b, ErrBranchExists
		}

		return b, fmt.Errorf("%w: %w", ErrCreateBranch, err)
	}

	b.New = true

	return b, nil
}

func (g GitHub) GetOrCreateBranch(ctx context.Context, repo Repo, branch, sha string) (Branch, error) {
	b, err := g.GetBranch(ctx, repo, branch)
	if err == nil {
		return b, nil
	}

	if !errors.Is(err, ErrNoBranch) {
		return b, err
	}

	return g.CreateBranch(ctx, repo, branch, sha)
}
