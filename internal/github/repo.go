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
func (g *GitHub) GetDefaultBranch(ctx context.Context, r Repo) (string, error) {
	if _, err := g.req(ctx, http.MethodGet, r.APIPath(), nil, &r); err != nil {
		return "", fmt.Errorf("%w: %w", errGetRepo, err)
	}

	return r.DefaultBranch, nil
}
