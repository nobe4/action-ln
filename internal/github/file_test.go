package github

import (
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

const (
	contentPath = "/repos/owner/repo/contents/path/to/file"
)

func TestGetFile(t *testing.T) {
	t.Parallel()

	repo := Repo{Owner: User{"owner"}, Repo: "repo"}

	t.Run("fails to decode the content", func(t *testing.T) {
		t.Parallel()

		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, `{"content": "_not base64"}`)
		}))

		g := New("token", ts.URL)

		_, err := g.GetFile(t.Context(), repo, "path/to/file")
		if !errors.Is(err, base64.CorruptInputError(0)) {
			t.Fatalf("expected base64 error, got %v", err)
		}
	})

	t.Run("succeeds", func(t *testing.T) {
		t.Parallel()

		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assertReq(t, r, http.MethodGet, contentPath, nil)

			fmt.Fprintln(w, `{"content": "b2s="}`)
		}))

		g := New("token", ts.URL)

		c, err := g.GetFile(t.Context(), repo, "path/to/file")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if c.Content != "ok" {
			t.Fatalf("expected content to be 'ok' but got %s", c.Content)
		}
	})
}

func TestUpdateFile(t *testing.T) {
	t.Parallel()

	repo := Repo{Owner: User{"owner"}, Repo: "repo"}
	content := File{
		Content: "ok",
		SHA:     "sha",
		Path:    "path/to/file",
	}

	t.Run("fails", func(t *testing.T) {
		t.Parallel()

		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))

		g := New("token", ts.URL)

		_, err := g.UpdateFile(t.Context(), repo, content, "branch", "message")
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("succeeds", func(t *testing.T) {
		t.Parallel()

		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assertReq(t, r,
				http.MethodPut,
				contentPath,
				[]byte(`{"message":"message","content":"b2s=","sha":"sha","branch":"branch"}`),
			)

			w.WriteHeader(http.StatusCreated)
			fmt.Fprintln(w, `{"content": {"sha":"newSha"}}`)
		}))

		g := New("token", ts.URL)

		c, err := g.UpdateFile(t.Context(), repo, content, "branch", "message")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if c.SHA != "newSha" {
			t.Fatalf("expected new sha to be 'newSha' but got '%s'", c.SHA)
		}

		if c.Content != content.Content {
			t.Fatalf("expected content to be '%s' but got '%s'", content.Content, c.Content)
		}
	})
}
