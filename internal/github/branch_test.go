package github

import (
	"fmt"
	"net/http"
	"testing"
)

const (
	branch        = "branch"
	sha           = "sha123"
	branchAPIPath = "/repos/owner/repo/branches/branch"
	refAPIPath    = "/repos/owner/repo/git/refs"
)

func TestGetBranch(t *testing.T) {
	t.Parallel()

	g := setup(t, func(w http.ResponseWriter, r *http.Request) {
		assertReq(t, r, http.MethodGet, branchAPIPath, nil)

		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"name": "%s", "commit": { "sha": "%s" } }`, branch, sha)
	})

	got, err := g.GetBranch(t.Context(), repo, branch)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if got.Name != branch {
		t.Fatalf("expected branch name to be '%s' but got '%s'", branch, got.Name)
	}

	if got.Commit.SHA != sha {
		t.Fatalf("expected commit SHA to be '%s' but got '%s'", sha, got.Commit.SHA)
	}
}

func TestCreateBranch(t *testing.T) {
	t.Parallel()

	t.Run("fails to create a new branch", func(t *testing.T) {
		t.Parallel()

		g := setup(t, func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusNotImplemented)
		})

		_, err := g.CreateBranch(t.Context(), repo, branch, sha)
		if err == nil {
			t.Fatalf("expected no error, got %v", err)
		}
	})

	t.Run("succeeds", func(t *testing.T) {
		t.Parallel()

		g := setup(t, func(w http.ResponseWriter, r *http.Request) {
			assertReq(t, r,
				http.MethodPost,
				refAPIPath,
				fmt.Appendf(nil, `{"ref":"refs/heads/%s","sha":"%s"}`, branch, sha),
			)

			w.WriteHeader(http.StatusCreated)
		})

		got, err := g.CreateBranch(t.Context(), repo, branch, sha)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if got.Name != branch || got.Commit.SHA != sha {
			t.Fatalf("want '%v', but got %v", branch, got)
		}
	})
}

func TestGetOrCreateBranch(t *testing.T) {
	t.Parallel()

	t.Run("finds existing branch", func(t *testing.T) {
		t.Parallel()

		g := setup(t, func(w http.ResponseWriter, r *http.Request) {
			assertReq(t, r, http.MethodGet, branchAPIPath, nil)
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, `{"name": "%s", "commit": { "sha": "%s" } }`, branch, sha)
		})

		got, err := g.GetOrCreateBranch(t.Context(), repo, branch, sha)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if got.Name != branch || got.Commit.SHA != sha {
			t.Fatalf("want '%v', but got %v", branch, got)
		}
	})

	t.Run("fails to get existing branch", func(t *testing.T) {
		t.Parallel()

		g := setup(t, func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		})

		_, err := g.GetOrCreateBranch(t.Context(), repo, branch, sha)
		if err == nil {
			t.Fatalf("expected error, got nil")
		}
	})

	t.Run("creates branch if it does not exist", func(t *testing.T) {
		t.Parallel()

		reqIndex := 0
		g := setup(t, func(w http.ResponseWriter, r *http.Request) {
			switch reqIndex {
			case 0:
				assertReq(t, r, http.MethodGet, branchAPIPath, nil)
				w.WriteHeader(http.StatusNotFound)
			case 1:
				assertReq(t, r, http.MethodPost, refAPIPath, nil)
				w.WriteHeader(http.StatusOK)
			}

			reqIndex++
		})

		got, err := g.GetOrCreateBranch(t.Context(), repo, branch, sha)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if got.Name != branch || got.Commit.SHA != sha {
			t.Fatalf("want '%v', but got %v", branch, got)
		}
	})
}
