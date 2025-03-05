package github

import (
	"fmt"
	"net/http"
	"testing"
)

func TestGetDefaultBranch(t *testing.T) {
	t.Parallel()

	g := setup(t, func(w http.ResponseWriter, r *http.Request) {
		assertReq(t, r, http.MethodGet, "/repos/owner/repo", nil)

		fmt.Fprintf(w, `{"default_branch": "%s"}`, branch)
	})

	b, err := g.GetDefaultBranch(t.Context(), repo)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if b != branch {
		t.Fatalf("expected default_branch to be '%s', got '%s'", branch, b)
	}

	if repo.DefaultBranch != "" {
		t.Fatalf("expected repo not to change its default branch, got %v", repo)
	}
}
