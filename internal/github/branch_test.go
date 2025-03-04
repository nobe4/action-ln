package github

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

const (
	branch        = "branch"
	sha           = "sha123"
	branchAPIPath = "/repos/owner/repo/branches/branch"
)

func TestGetBranch(t *testing.T) {
	t.Parallel()

	repo := Repo{Owner: User{Login: "owner"}, Repo: "repo"}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertReq(t, r, http.MethodGet, branchAPIPath, nil)

		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, `{"name": "branch", "commit": { "sha": "sha123" } }`)
	}))

	g := New("token", ts.URL)

	got, err := g.GetBranch(t.Context(), repo, branch)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if got.Name != branch {
		t.Fatalf("expected branch name to be '%s' but got %s", branch, got.Name)
	}

	if got.Commit.SHA != "sha123" {
		t.Fatalf("expected commit SHA to be 'sha123' but got %s", got.Commit.SHA)
	}
}

func TestCreateBranch(t *testing.T) {
	t.Parallel()

	repo := Repo{Owner: User{Login: "owner"}, Repo: "repo"}

	t.Run("fails to create a new branch", func(t *testing.T) {
		t.Parallel()

		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotImplemented)
		}))

		g := New("token", ts.URL)

		_, err := g.CreateBranch(t.Context(), repo, branch, sha)
		if err == nil {
			t.Fatalf("expected no error, got %v", err)
		}
	})

	t.Run("succeeds", func(t *testing.T) {
		t.Parallel()

		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assertReq(t, r,
				http.MethodPost,
				"/repos/owner/repo/git/refs",
				[]byte(`{"ref":"refs/heads/branch","sha":"sha123"}`),
			)

			w.WriteHeader(http.StatusCreated)
		}))

		g := New("token", ts.URL)

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

	repo := Repo{Owner: User{Login: "owner"}, Repo: "repo"}

	t.Run("finds existing branch", func(t *testing.T) {
		t.Parallel()

		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assertReq(t, r, http.MethodGet, "/repos/owner/repo/branches/branch", nil)
			w.WriteHeader(http.StatusOK)
			fmt.Fprintln(w, `{"name": "branch", "commit": { "sha": "sha123" } }`)
		}))

		g := New("token", ts.URL)

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

		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))

		g := New("token", ts.URL)

		_, err := g.GetOrCreateBranch(t.Context(), repo, branch, sha)
		if err == nil {
			t.Fatalf("expected error, got nil")
		}
	})

	t.Run("creates branch if it does not exist", func(t *testing.T) {
		t.Parallel()

		reqIndex := 0
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch reqIndex {
			case 0:
				assertReq(t, r, http.MethodGet, "/repos/owner/repo/branches/branch", nil)
				w.WriteHeader(http.StatusNotFound)
			case 1:
				assertReq(t, r, http.MethodPost, "/repos/owner/repo/git/refs", nil)
				w.WriteHeader(http.StatusOK)
			}

			reqIndex++
		}))

		g := New("token", ts.URL)

		got, err := g.GetOrCreateBranch(t.Context(), repo, branch, sha)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if got.Name != branch || got.Commit.SHA != sha {
			t.Fatalf("want '%v', but got %v", branch, got)
		}
	})
}
