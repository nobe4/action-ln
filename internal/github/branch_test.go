package github

import (
	"errors"
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

	t.Run("missing branch", func(t *testing.T) {
		t.Parallel()

		g := setup(t, func(w http.ResponseWriter, r *http.Request) {
			assertReq(t, r, http.MethodGet, branchAPIPath, nil)
			w.WriteHeader(http.StatusNotFound)
		})

		_, err := g.GetBranch(t.Context(), repo, branch)
		if !errors.Is(err, ErrNoBranch) {
			t.Fatalf("expected error %v, got %v", ErrNoBranch, err)
		}
	})

	t.Run("server error", func(t *testing.T) {
		t.Parallel()

		g := setup(t, func(w http.ResponseWriter, r *http.Request) {
			assertReq(t, r, http.MethodGet, branchAPIPath, nil)
			w.WriteHeader(http.StatusInternalServerError)
		})

		_, err := g.GetBranch(t.Context(), repo, branch)
		if !errors.Is(err, ErrGetBranch) {
			t.Fatalf("expected error %v, got %v", ErrNoBranch, err)
		}
	})

	t.Run("succeeds", func(t *testing.T) {
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
	})
}

func TestCreateBranch(t *testing.T) {
	t.Parallel()

	t.Run("branch exists", func(t *testing.T) {
		t.Parallel()

		g := setup(t, func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusUnprocessableEntity)
		})

		_, err := g.CreateBranch(t.Context(), repo, branch, sha)
		if !errors.Is(err, ErrBranchExists) {
			t.Fatalf("expected error %v, got %v", ErrCreateBranch, err)
		}
	})

	t.Run("server error", func(t *testing.T) {
		t.Parallel()

		g := setup(t, func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		})

		_, err := g.CreateBranch(t.Context(), repo, branch, sha)
		if !errors.Is(err, ErrCreateBranch) {
			t.Fatalf("expected error %v, got %v", ErrCreateBranch, err)
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

func TestDeleteBranch(t *testing.T) {
	t.Parallel()

	refPath := refAPIPath + "/heads/" + branch

	t.Run("server error", func(t *testing.T) {
		t.Parallel()

		g := setup(t, func(w http.ResponseWriter, r *http.Request) {
			assertReq(t, r, http.MethodDelete, refPath, nil)
			w.WriteHeader(http.StatusInternalServerError)
		})

		err := g.DeleteBranch(t.Context(), repo, branch)
		if !errors.Is(err, ErrDeleteBranch) {
			t.Fatalf("expected error %v, got %v", ErrNoBranch, err)
		}
	})

	t.Run("succeeds", func(t *testing.T) {
		t.Parallel()
		g := setup(t, func(w http.ResponseWriter, r *http.Request) {
			assertReq(t, r, http.MethodDelete, refPath, nil)

			w.WriteHeader(http.StatusOK)
		})

		err := g.DeleteBranch(t.Context(), repo, branch)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
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

		i := 0
		g := setup(t, func(w http.ResponseWriter, r *http.Request) {
			switch i {
			case 0:
				assertReq(t, r, http.MethodGet, branchAPIPath, nil)
				w.WriteHeader(http.StatusNotFound)
			case 1:
				assertReq(t, r, http.MethodPost, refAPIPath, nil)
				w.WriteHeader(http.StatusOK)
			default:
				w.WriteHeader(http.StatusInternalServerError)
			}

			i++
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
