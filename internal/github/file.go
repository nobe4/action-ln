package github

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
)

type File struct {
	Name       string `json:"name"`
	Path       string `json:"path"`
	RawContent string `json:"content"`
	SHA        string `json:"sha"`
	Content    string
}

// https://docs.github.com/en/rest/repos/contents?apiVersion=2022-11-28#get-repository-content
func (g GitHub) GetFile(ctx context.Context, repo Repo, path string) (File, error) {
	path = fmt.Sprintf("/repos/%s/%s/contents/%s", repo.Owner.Login, repo.Repo, path)

	c := File{}
	if _, err := g.req(ctx, http.MethodGet, path, nil, &c); err != nil {
		// TODO: make constant error
		return File{}, fmt.Errorf("failed to get file: %w", err)
	}

	decoded, err := base64.StdEncoding.DecodeString(c.RawContent)
	if err != nil {
		// TODO: make constant error
		return File{}, fmt.Errorf("failed to decode content: %w", err)
	}

	c.Content = string(decoded)

	return c, nil
}

// https://docs.github.com/en/rest/repos/contents?apiVersion=2022-11-28#create-or-update-file-contents
func (g GitHub) UpdateFile(ctx context.Context, repo Repo, f File, branch, message string) (File, error) {
	body, err := json.Marshal(struct {
		Message string `json:"message"`
		Content string `json:"content"`
		SHA     string `json:"sha"`
		Branch  string `json:"branch"`
	}{
		Message: message,
		Content: base64.StdEncoding.EncodeToString([]byte(f.Content)),
		Branch:  branch,
		SHA:     f.SHA,
	})
	if err != nil {
		// TODO: make constant error
		return File{}, fmt.Errorf("failed to marshal request: %w", err)
	}

	// NOTE: The response wrapes the content in an extra `{ "content": {} }`.
	// It's not technically the same as the File we get from GetFile. For
	// this purpose it's enough, we can reinject the `Content` value after.
	out := struct {
		File File `json:"content"`
	}{}

	if _, err := g.req(
		ctx,
		http.MethodPut,
		fmt.Sprintf("/repos/%s/%s/contents/%s", repo.Owner.Login, repo.Repo, f.Path),
		bytes.NewReader(body),
		&out,
	); err != nil {
		// TODO: make constant error
		return File{}, fmt.Errorf("failed to create branch: %w", err)
	}

	out.File.Content = f.Content

	return out.File, nil
}
