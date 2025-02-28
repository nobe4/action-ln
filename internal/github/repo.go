package github

import (
	"context"
	"encoding/base64"
	"fmt"
)

type Repo struct {
	Owner         string `json:"owner"`
	Repo          string `json:"repo"`
	DefaultBranch string `json:"default_branch"`
}

// https://docs.github.com/en/rest/repos/repos?apiVersion=2022-11-28
func (g GitHub) GetDefaultBranch(ctx context.Context, repo Repo) (string, error) {
	path := fmt.Sprintf("/repos/%s/%s", repo.Owner, repo.Repo)

	if err := g.req(ctx, "GET", path, nil, &repo); err != nil {
		return "", fmt.Errorf("failed to get repo: %w", err)
	}

	return repo.DefaultBranch, nil
}

type Content struct {
	Name       string `json:"name"`
	Path       string `json:"path"`
	RawContent string `json:"content"`
	Content    string
}

// https://docs.github.com/en/rest/repos/contents?apiVersion=2022-11-28
func (g GitHub) GetContent(ctx context.Context, repo Repo, path string) (Content, error) {
	c := Content{}

	path = fmt.Sprintf("/repos/%s/%s/contents/%s", repo.Owner, repo.Repo, path)

	if err := g.req(ctx, "GET", path, nil, &c); err != nil {
		return c, fmt.Errorf("failed to get user: %w", err)
	}

	decoded, err := base64.StdEncoding.DecodeString(c.RawContent)
	if err != nil {
		return c, fmt.Errorf("failed to decode content: %w", err)
	}

	c.Content = string(decoded)

	return c, nil
}
