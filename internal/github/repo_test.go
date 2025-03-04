package github

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetDefaultBranch(t *testing.T) {
	t.Parallel()

	repo := Repo{Owner: User{Login: "owner"}, Repo: "repo"}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/repos/owner/repo" {
			t.Fatal("invalid path", r.URL.Path)
		}

		if r.Method != http.MethodGet {
			t.Fatal("invalid method", r.Method)
		}

		fmt.Fprintln(w, `{"default_branch": "main"}`)
	}))

	g := New("token", ts.URL)

	b, err := g.GetDefaultBranch(t.Context(), repo)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if b != "main" {
		t.Fatalf("expected default_branch to be 'main' but got %s", b)
	}

	// ensures that this doesn't modify the repo passed as an argument.
	if repo.DefaultBranch != "" {
		t.Fatalf("expected repo not to change got %v", repo)
	}
}
