package github

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetBranch(t *testing.T) {
	t.Parallel()

	repo := Repo{Owner: User{Login: "owner"}, Repo: "repo"}
	branch := "branch"

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/repos/owner/repo/branches/branch" {
			t.Fatal("invalid path", r.URL.Path)
		}

		fmt.Fprintln(w, `{"name": "branch", "commit": { "sha": "123abc" } }`)
	}))

	g := New("token", ts.URL)

	got, err := g.GetBranch(t.Context(), repo, branch)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if got.Name != branch {
		t.Fatalf("expected branch name to be '%s' but got %s", branch, got.Name)
	}

	if got.Commit.SHA != "123abc" {
		t.Fatalf("expected commit SHA to be '123abc' but got %s", got.Commit.SHA)
	}
}
