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
	ErrDeleteBranch = errors.New("failed to delete branch")
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
func (g *GitHub) GetBranch(ctx context.Context, r Repo, name string) (Branch, error) {
	b := Branch{}

	path := fmt.Sprintf("/repos/%s/branches/%s", r, name)

	if status, err := g.req(ctx, http.MethodGet, path, nil, &b); err != nil {
		if status == http.StatusNotFound {
			return Branch{}, ErrNoBranch
		}

		return Branch{}, fmt.Errorf("%w: %w", ErrGetBranch, err)
	}

	b.New = false

	return b, nil
}

// https://docs.github.com/en/rest/git/refs?apiVersion=2022-11-28#create-a-reference
func (g *GitHub) CreateBranch(ctx context.Context, r Repo, branch, sha string) (Branch, error) {
	b := Branch{
		Name: branch,
		Commit: Commit{
			SHA: sha,
		},
	}

	path := fmt.Sprintf("/repos/%s/git/refs", r)

	body, err := json.Marshal(struct {
		Ref string `json:"ref"`
		SHA string `json:"sha"`
	}{
		Ref: "refs/heads/" + branch,
		SHA: sha,
	})
	if err != nil {
		return Branch{}, fmt.Errorf("%w: %w", ErrMarshalRequest, err)
	}

	if status, err := g.req(ctx, http.MethodPost, path, bytes.NewReader(body), nil); err != nil {
		if status == http.StatusUnprocessableEntity {
			return Branch{}, ErrBranchExists
		}

		return Branch{}, fmt.Errorf("%w: %w", ErrCreateBranch, err)
	}

	b.New = true

	return b, nil
}

// https://docs.github.com/en/rest/git/refs?apiVersion=2022-11-28#delete-a-reference
func (g *GitHub) DeleteBranch(ctx context.Context, r Repo, name string) error {
	path := fmt.Sprintf("/repos/%s/git/refs/heads/%s", r, name)

	if _, err := g.req(ctx, http.MethodDelete, path, nil, nil); err != nil {
		return fmt.Errorf("%w: %w", ErrDeleteBranch, err)
	}

	return nil
}

func (g *GitHub) GetOrCreateBranch(ctx context.Context, r Repo, name, sha string) (Branch, error) {
	b, err := g.GetBranch(ctx, r, name)
	if err == nil {
		return b, nil
	}

	if !errors.Is(err, ErrNoBranch) {
		return b, err
	}

	return g.CreateBranch(ctx, r, name, sha)
}
