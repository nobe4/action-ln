package github

import (
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetContent(t *testing.T) {
	t.Parallel()

	repo := Repo{Owner: User{"owner"}, Repo: "repo"}

	t.Run("fails to decode the content", func(t *testing.T) {
		t.Parallel()

		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/repos/owner/repo/contents/path/to/file" {
				t.Fatal("invalid path", r.URL.Path)
			}

			fmt.Fprintln(w, `{"content": "_not base64"}`)
		}))

		g := New("token", ts.URL)

		_, err := g.GetContent(t.Context(), repo, "path/to/file")
		if !errors.Is(err, base64.CorruptInputError(0)) {
			t.Fatalf("expected base64 error, got %v", err)
		}
	})

	t.Run("succeeds", func(t *testing.T) {
		t.Parallel()

		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/repos/owner/repo/contents/path/to/file" {
				t.Fatal("invalid path", r.URL.Path)
			}

			fmt.Fprintln(w, `{"content": "b2s="}`)
		}))

		g := New("token", ts.URL)

		c, err := g.GetContent(t.Context(), repo, "path/to/file")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if c.Content != "ok" {
			t.Fatalf("expected content to be 'ok' but got %s", c.Content)
		}
	})
}

func TestCreateOrUpdateContent(t *testing.T) {
	t.Parallel()

	repo := Repo{Owner: User{"owner"}, Repo: "repo"}
	content := Content{
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

		_, err := g.CreateOrUpdateContent(t.Context(), repo, content, "branch", "message")
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("succeeds", func(t *testing.T) {
		t.Parallel()

		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/repos/owner/repo/contents/path/to/file" {
				t.Fatal("invalid path", r.URL.Path)
			}

			if r.Method != http.MethodPut {
				t.Fatal("invalid method", r.Method)
			}

			body, err := io.ReadAll(r.Body)
			if err != nil {
				t.Fatal("failed to read body", err)
			}

			if string(body) != `{"message":"message","content":"b2s=","sha":"sha","branch":"branch"}` {
				t.Fatal("invalid body", string(body))
			}

			w.WriteHeader(http.StatusCreated)
			fmt.Fprintln(w, `{"content": {"sha":"newSha"}}`)
		}))

		g := New("token", ts.URL)

		c, err := g.CreateOrUpdateContent(t.Context(), repo, content, "branch", "message")
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
