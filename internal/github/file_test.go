package github

import (
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"testing"
)

const (
	filePath      = "path/to/file"
	contentPath   = "/repos/owner/repo/contents/" + filePath
	content       = "ok"
	base64Content = "b2s="
	message       = "message"
)

func TestGetFile(t *testing.T) {
	t.Parallel()

	t.Run("fails to decode the content", func(t *testing.T) {
		t.Parallel()

		g := setup(t, func(w http.ResponseWriter, _ *http.Request) {
			fmt.Fprintln(w, `{"content": "_not base64"}`)
		})

		f := File{Repo: repo, Path: filePath}
		if err := g.GetFile(t.Context(), &f); !errors.Is(err, base64.CorruptInputError(0)) {
			t.Fatalf("expected base64 error, got %v", err)
		}
	})

	t.Run("succeeds", func(t *testing.T) {
		t.Parallel()

		g := setup(t, func(w http.ResponseWriter, r *http.Request) {
			assertReq(t, r, http.MethodGet, contentPath, nil)

			fmt.Fprintf(w, `{"content": "%s"}`, base64Content)
		})

		f := File{Repo: repo, Path: filePath}

		if err := g.GetFile(t.Context(), &f); err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if f.Content != "ok" {
			t.Fatalf("expected content to be 'ok' but got %s", f.Content)
		}
	})
}

func TestUpdateFile(t *testing.T) {
	t.Parallel()

	content := File{
		Content: content,
		SHA:     sha,
		Path:    filePath,
	}

	t.Run("fails", func(t *testing.T) {
		t.Parallel()

		g := setup(t, func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		})

		_, err := g.UpdateFile(t.Context(), repo, content, branch, message)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("succeeds", func(t *testing.T) {
		t.Parallel()

		const newSha = "newSha"

		g := setup(t, func(w http.ResponseWriter, r *http.Request) {
			assertReq(t, r,
				http.MethodPut,
				contentPath,
				fmt.Appendf(nil, `{"message":"%s","content":"%s","sha":"%s","branch":"%s"}`, message, base64Content, sha, branch),
			)

			w.WriteHeader(http.StatusCreated)
			fmt.Fprintf(w, `{"content": {"sha":"%s"}}`, newSha)
		})

		c, err := g.UpdateFile(t.Context(), repo, content, branch, message)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if c.SHA != newSha {
			t.Fatalf("expected new sha to be '%s' but got '%s'", newSha, c.SHA)
		}

		if c.Content != content.Content {
			t.Fatalf("expected content to be '%s' but got '%s'", content.Content, c.Content)
		}
	})
}
