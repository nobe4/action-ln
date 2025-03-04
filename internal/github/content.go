package github

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
)

type Content struct {
	Name       string `json:"name"`
	Path       string `json:"path"`
	RawContent string `json:"content"`
	SHA        string `json:"sha"`
	Content    string
}

// https://docs.github.com/en/rest/repos/contents?apiVersion=2022-11-28#get-repository-content
func (g GitHub) GetContent(ctx context.Context, repo Repo, path string) (Content, error) {
	path = fmt.Sprintf("/repos/%s/%s/contents/%s", repo.Owner.Login, repo.Repo, path)

	c := Content{}
	if _, err := g.req(ctx, http.MethodGet, path, nil, &c); err != nil {
		return Content{}, fmt.Errorf("failed to get user: %w", err)
	}

	decoded, err := base64.StdEncoding.DecodeString(c.RawContent)
	if err != nil {
		return Content{}, fmt.Errorf("failed to decode content: %w", err)
	}

	c.Content = string(decoded)

	return c, nil
}

// https://docs.github.com/en/rest/repos/contents?apiVersion=2022-11-28#create-or-update-file-contents
func (g GitHub) CreateOrUpdateContent(
	ctx context.Context,
	repo Repo,
	c Content,
	branch,
	message string,
) (Content, error) {
	body, err := json.Marshal(struct {
		Message string `json:"message"`
		Content string `json:"content"`
		SHA     string `json:"sha"`
		Branch  string `json:"branch"`
	}{
		Message: message,
		Content: base64.StdEncoding.EncodeToString([]byte(c.Content)),
		Branch:  branch,
		SHA:     c.SHA,
	})
	if err != nil {
		return Content{}, fmt.Errorf("failed to marshal request: %w", err)
	}

	// The response wrapes the content in an extra `{ "content": {} }`.
	out := struct {
		Content Content `json:"content"`
	}{}

	if _, err := g.req(
		ctx,
		http.MethodPut,
		fmt.Sprintf("/repos/%s/%s/contents/%s", repo.Owner.Login, repo.Repo, c.Path),
		bytes.NewReader(body),
		&out,
	); err != nil {
		return Content{}, fmt.Errorf("failed to create branch: %w", err)
	}

	return out.Content, nil
}
