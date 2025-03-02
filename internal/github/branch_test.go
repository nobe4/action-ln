package github

import (
	"bytes"
	"fmt"
	"io"
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

func TestCreateBranch(t *testing.T) {
	t.Parallel()

	repo := Repo{Owner: User{Login: "owner"}, Repo: "repo"}
	branch := "branch"
	sha := "123abc"

	t.Run("fails to create a new branch", func(t *testing.T) {
		t.Parallel()

		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/repos/owner/repo/git/refs" {
				t.Fatal("invalid path", r.URL.Path)
			}

			if r.Method != http.MethodPost {
				t.Fatal("invalid method", r.Method)
			}

			body, err := io.ReadAll(r.Body)
			if err != nil {
				t.Fatal("failed to read body", err)
			}

			if !bytes.Equal(body, []byte(`{"ref":"refs/heads/branch","sha":"123abc"}`)) {
				t.Fatal("invalid body", string(body))
			}

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
