package github

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/nobe4/action-ln/internal/log"
)

type Repo struct {
	Owner         User   `json:"owner"`
	Repo          string `json:"repo"`
	DefaultBranch string `json:"default_branch"`
}

var errGetRepo = errors.New("failed to get repo")

func (r Repo) Equal(o Repo) bool {
	return r.Repo == o.Repo && r.Owner.Login == o.Owner.Login
}

func (r Repo) Empty() bool {
	return r.Repo == "" && r.Owner.Login == "" && r.DefaultBranch == ""
}

func (r Repo) String() string {
	return fmt.Sprintf("%s/%s", r.Owner.Login, r.Repo)
}

func (r Repo) APIPath() string {
	return fmt.Sprintf("/repos/%s", r)
}

// https://docs.github.com/en/rest/repos/repos?apiVersion=2022-11-28
func (g *GitHub) GetDefaultBranchName(ctx context.Context, r Repo) (string, error) {
	log.Debug("Get default branch name", "repo", r)

	if _, err := g.req(ctx, http.MethodGet, r.APIPath(), nil, &r); err != nil {
		return "", fmt.Errorf("%w: %w", errGetRepo, err)
	}

	return r.DefaultBranch, nil
}

func (g *GitHub) GetDefaultBranch(ctx context.Context, r Repo) (Branch, error) {
	log.Debug("Get default branch", "repo", r)

	name, err := g.GetDefaultBranchName(ctx, r)
	if err != nil {
		return Branch{}, err
	}

	b, err := g.GetBranch(ctx, r, name)
	if err != nil {
		return Branch{}, err
	}

	return b, nil
}

func (g *GitHub) GetBaseAndHeadBranches(ctx context.Context, r Repo, headName string) (
	base Branch, head Branch,
	err error,
) {
	if base, err = g.GetDefaultBranch(ctx, r); err != nil {
		return base, head, err
	}

	if head, err = g.GetOrCreateBranch(ctx, r, headName, base.Commit.SHA); err != nil {
		return base, head, err
	}

	return base, head, nil
}
