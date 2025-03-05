package github

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

const (
	head         = "head"
	base         = "base"
	title        = "title"
	body         = "body"
	number       = 123
	pullAPIPath  = "/repos/owner/repo/pulls"
	pullAPIQuery = "base=base&head=owner%3Ahead&per_page=1&state=open"
)

var repo = Repo{Owner: User{Login: "owner"}, Repo: "repo"}

func TestGetPull(t *testing.T) {
	t.Parallel()

	t.Run("finds a pull", func(t *testing.T) {
		t.Parallel()

		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assertReq(t, r, http.MethodGet, pullAPIPath, nil)

			if r.URL.RawQuery != pullAPIQuery {
				t.Fatalf("expected query to be '%s' but got '%s'", pullAPIQuery, r.URL.RawQuery)
			}

			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, `[{"number": %d}]\n`, number)
		}))

		g := New("token", ts.URL)

		got, err := g.GetPull(t.Context(), repo, base, head)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if got.Number != number {
			t.Fatalf("expected number to be %d but got %d", number, got.Number)
		}
	})

	t.Run("finds no pull", func(t *testing.T) {
		t.Parallel()

		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			fmt.Fprintln(w, `[]`)
		}))

		g := New("token", ts.URL)

		_, err := g.GetPull(t.Context(), repo, "base", "head")
		if !errors.Is(err, errNoPull) {
			t.Fatalf("expected error to be %q, got %q", errNoPull, err)
		}
	})
}

func TestCreatePull(t *testing.T) {
	t.Parallel()

	repo := Repo{Owner: User{Login: "owner"}, Repo: "repo"}

	t.Run("pull exists", func(t *testing.T) {
		t.Parallel()

		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusUnprocessableEntity)
		}))

		g := New("token", ts.URL)

		_, err := g.CreatePull(t.Context(), repo, "base", "head", "title", "body")
		if !errors.Is(err, errPullExists) {
			t.Fatalf("expected error %q, got %q", errPullExists, err)
		}
	})

	t.Run("create a pull", func(t *testing.T) {
		t.Parallel()

		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assertReq(t, r,
				http.MethodPost,
				pullAPIPath,
				[]byte(`{"title":"title","head":"owner:head","base":"base","body":"body"}`),
			)

			w.WriteHeader(http.StatusCreated)
			fmt.Fprintf(w, `{"number": %d}\n`, number)
		}))

		g := New("token", ts.URL)

		got, err := g.CreatePull(t.Context(), repo, base, head, title, body)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if got.Number != number {
			t.Fatalf("expected number to be %d but got %d", number, got.Number)
		}
	})
}

func TestGetOrCreatePull(t *testing.T) {
	t.Parallel()

	t.Run("finds existing pull", func(t *testing.T) {
		t.Parallel()

		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.RawQuery != pullAPIQuery {
				t.Fatalf("expected query to be '%s' but got '%s'", pullAPIQuery, r.URL.RawQuery)
			}

			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, `[{"number": %d}]\n`, number)
		}))

		g := New("token", ts.URL)

		got, err := g.GetOrCreatePull(t.Context(), repo, base, head, title, body)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if got.Number != number {
			t.Fatalf("want number %d, but got %d", number, got.Number)
		}
	})

	t.Run("fails to get existing pull", func(t *testing.T) {
		t.Parallel()

		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))

		g := New("token", ts.URL)

		_, err := g.GetOrCreatePull(t.Context(), repo, base, head, title, body)
		if err == nil {
			t.Fatalf("expected error, got nil")
		}
	})

	t.Run("creates pull if it does not exist", func(t *testing.T) {
		t.Parallel()

		reqIndex := 0
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch reqIndex {
			case 0:
				assertReq(t, r, http.MethodGet, pullAPIPath, nil)
				w.WriteHeader(http.StatusOK)
				fmt.Fprintln(w, `[]`)
			case 1:
				assertReq(t, r, http.MethodPost, pullAPIPath, nil)
				w.WriteHeader(http.StatusOK)
				fmt.Fprintf(w, `{"number": %d}\n`, number)
			}

			reqIndex++
		}))

		g := New("token", ts.URL)

		got, err := g.GetOrCreatePull(t.Context(), repo, base, head, title, body)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if got.Number != number {
			t.Fatalf("expected number to be %d but got %d", number, got.Number)
		}
	})
}
