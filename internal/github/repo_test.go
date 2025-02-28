package github

import (
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetContent(t *testing.T) {
	t.Parallel()

	repo := Repo{Owner: "owner", Repo: "repo"}

	t.Run("fails to decode the content", func(t *testing.T) {
		t.Parallel()

		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
