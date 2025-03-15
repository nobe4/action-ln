package github

import (
	"fmt"
	"net/http"
	"testing"
)

func TestGetDefaultBranchName(t *testing.T) {
	t.Parallel()

	g := setup(t, func(w http.ResponseWriter, r *http.Request) {
		assertReq(t, r, http.MethodGet, "/repos/owner/repo", nil)

		fmt.Fprintf(w, `{"default_branch": "%s"}`, branch)
	})

	b, err := g.GetDefaultBranchName(t.Context(), repo)
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

func TestGetDefaultBranch(t *testing.T) {
	t.Parallel()

	i := 0
	g := setup(t, func(w http.ResponseWriter, r *http.Request) {
		switch i {
		case 0:
			assertReq(t, r, http.MethodGet, "/repos/owner/repo", nil)
			fmt.Fprintf(w, `{"default_branch": "%s"}`, branch)
		case 1:
			assertReq(t, r, http.MethodGet, "/repos/owner/repo/branches/"+branch, nil)
			fmt.Fprintf(w, `{"name": "%s","commit":{"sha":"%s"}}`, branch, sha)
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}

		i++
	})

	b, err := g.GetDefaultBranch(t.Context(), repo)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if b.Name != branch {
		t.Fatalf("expected branch to be '%s', got '%s'", branch, b.Name)
	}

	if b.Commit.SHA != sha {
		t.Fatalf("expected sha to be '%s', got '%s'", sha, b.Commit.SHA)
	}
}

func TestGetBaseAndHeadBranches(t *testing.T) {
	t.Parallel()

	baseName := "base"
	headName := "head"

	i := 0
	g := setup(t, func(w http.ResponseWriter, r *http.Request) {
		switch i {
		case 0:
			assertReq(t, r, http.MethodGet, "/repos/owner/repo", nil)
			fmt.Fprintf(w, `{"default_branch": "%s"}`, baseName)
		case 1:
			assertReq(t, r, http.MethodGet, "/repos/owner/repo/branches/"+baseName, nil)
			fmt.Fprintf(w, `{"name": "%s","commit":{"sha":"%s"}}`, baseName, sha)
		case 2:
			assertReq(t, r, http.MethodGet, "/repos/owner/repo/branches/"+headName, nil)
			w.WriteHeader(http.StatusNotFound)
		case 3:
			assertReq(t, r,
				http.MethodPost,
				"/repos/owner/repo/git/refs",
				fmt.Appendf(nil, `{"ref":"refs/heads/%s","sha":"%s"}`, headName, sha),
			)
			w.WriteHeader(http.StatusOK)
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}

		i++
	})

	base, head, err := g.GetBaseAndHeadBranches(t.Context(), repo, headName)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if base.Name != baseName {
		t.Fatalf("expected branch to be '%s', got '%s'", baseName, base.Name)
	}

	if base.New {
		t.Fatalf("expected branch not to be new, got %v", base.New)
	}

	if head.Name != headName {
		t.Fatalf("expected branch to be '%s', got '%s'", headName, head.Name)
	}

	if !head.New {
		t.Fatalf("expected branch to be new, got %v", head.New)
	}
}
