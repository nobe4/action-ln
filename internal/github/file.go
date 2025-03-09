package github

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type File struct {
	// Content from the API
	Name    string `json:"name"`
	Path    string `json:"path"`
	Content string `json:"content"`
	SHA     string `json:"sha"`

	// Content from the config
	Repo   Repo   `json:"repo"`
	Ref    string `json:"ref"`
	Commit string `json:"commit"`
}

var (
	ErrGetFile    = errors.New("failed to get file")
	ErrUpdateFile = errors.New("failed to create/update file")
	ErrDecodeFile = errors.New("failed to decode file")
)

func (f File) Equal(o File) bool {
	return f.Repo.Equal(o.Repo) && f.Path == o.Path && f.SHA == o.SHA && f.Commit == o.Commit
}

// https://docs.github.com/en/rest/repos/contents?apiVersion=2022-11-28#get-repository-content
// TODO: pass a file instead of repo + path.
func (g *GitHub) GetFile(ctx context.Context, repo Repo, path string) (File, error) {
	path = fmt.Sprintf("/repos/%s/%s/contents/%s", repo.Owner.Login, repo.Repo, path)

	c := File{}
	if _, err := g.req(ctx, http.MethodGet, path, nil, &c); err != nil {
		return File{}, fmt.Errorf("%w: %w", ErrGetFile, err)
	}

	decoded, err := base64.StdEncoding.DecodeString(c.Content)
	if err != nil {
		return File{}, fmt.Errorf("%w: %w", ErrDecodeFile, err)
	}

	c.Content = string(decoded)

	return c, nil
}

// https://docs.github.com/en/rest/repos/contents?apiVersion=2022-11-28#create-or-update-file-contents
// TODO: pass a file instead of repo + file.
// TODO: check that we're creating a new file correctly with the updated commit
// and sha.
func (g *GitHub) UpdateFile(ctx context.Context, repo Repo, f File, branch, message string) (File, error) {
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
		return File{}, fmt.Errorf("%w: %w", ErrMarshalRequest, err)
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
		return File{}, fmt.Errorf("%w: %w", ErrUpdateFile, err)
	}

	out.File.Content = f.Content

	return out.File, nil
}
