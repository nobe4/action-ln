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
	SHA     string `json:"sha"` // Blob hash.

	// Content from the config
	Repo   Repo   `json:"repo"`
	Ref    string `json:"ref"`
	Commit string `json:"commit"` // Commit hash.
}

var (
	ErrGetFile     = errors.New("failed to get file")
	ErrMissingFile = errors.New("file does not exist")
	ErrUpdateFile  = errors.New("failed to create/update file")
	ErrDecodeFile  = errors.New("failed to decode file")
)

func (f File) String() string {
	return fmt.Sprintf("%s:%s@%s", f.Repo, f.Path, f.Ref)
}

func (f File) Equal(o File) bool {
	return f.Repo.Equal(o.Repo) && f.Path == o.Path && f.SHA == o.SHA && f.Commit == o.Commit
}

func (f File) APIPath() string {
	return fmt.Sprintf("/repos/%s/contents/%s", f.Repo, f.Path)
}

// https://docs.github.com/en/rest/repos/contents?apiVersion=2022-11-28#get-repository-content
func (g *GitHub) GetFile(ctx context.Context, f *File) error {
	if status, err := g.req(ctx,
		http.MethodGet,
		f.APIPath(),
		nil,
		&f,
	); err != nil {
		if status == http.StatusNotFound {
			return fmt.Errorf("%w: %w", ErrMissingFile, err)
		}

		return fmt.Errorf("%w: %w", ErrGetFile, err)
	}

	decoded, err := base64.StdEncoding.DecodeString(f.Content)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrDecodeFile, err)
	}

	f.Content = string(decoded)

	return nil
}

// https://docs.github.com/en/rest/repos/contents?apiVersion=2022-11-28#create-or-update-file-contents
func (g *GitHub) UpdateFile(ctx context.Context, f File, branch, message string) (File, error) {
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

	// NOTE: Non-trivial update.
	// The response for this call will update the file directly. It's fine
	// because we want the new `SHA`. The `Name`, and `Path` won't change
	// because we only update the content of the file. Also, since we're passing
	// a file value, the original file won't get changed; we are creating a new
	// file.
	out := struct {
		File File `json:"content"`
	}{File: f}

	if _, err := g.req(
		ctx,
		http.MethodPut,
		f.APIPath(),
		bytes.NewReader(body),
		&out,
	); err != nil {
		return File{}, fmt.Errorf("%w: %w", ErrUpdateFile, err)
	}

	return out.File, nil
}
