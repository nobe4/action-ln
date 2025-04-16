package github

import (
	"errors"
	"fmt"
	"net/http"
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

func TestGetPull(t *testing.T) {
	t.Parallel()

	t.Run("finds a pull", func(t *testing.T) {
		t.Parallel()

		g := setup(t, func(w http.ResponseWriter, r *http.Request) {
			assertReq(t, r, http.MethodGet, pullAPIPath, nil)

			if r.URL.RawQuery != pullAPIQuery {
				t.Fatalf("expected query to be '%s' but got '%s'", pullAPIQuery, r.URL.RawQuery)
			}

			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, `[{"number": %d}]\n`, number)
		})

		got, err := g.GetPull(t.Context(), repo, base, head)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if got.Number != number {
			t.Fatalf("expected number to be %d but got %d", number, got.Number)
		}

		if got.Repo != repo {
			t.Fatalf("expected repo to be %v but got %v", repo, got.Repo)
		}

		if got.New {
			t.Fatalf("expected pull to be not new, but it is")
		}
	})

	t.Run("finds no pull", func(t *testing.T) {
		t.Parallel()

		g := setup(t, func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusOK)
			fmt.Fprintln(w, `[]`)
		})

		_, err := g.GetPull(t.Context(), repo, "base", "head")
		if !errors.Is(err, ErrNoPull) {
			t.Fatalf("expected error to be %q, got %q", ErrNoPull, err)
		}
	})
}

func TestCreatePull(t *testing.T) {
	t.Parallel()

	t.Run("pull exists", func(t *testing.T) {
		t.Parallel()

		g := setup(t, func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusUnprocessableEntity)
		})

		_, err := g.CreatePull(t.Context(), repo, "base", "head", "title", "body")
		if !errors.Is(err, ErrPullExists) {
			t.Fatalf("expected error %q, got %q", ErrPullExists, err)
		}
	})

	t.Run("create a pull", func(t *testing.T) {
		t.Parallel()

		g := setup(t, func(w http.ResponseWriter, r *http.Request) {
			assertReq(t, r,
				http.MethodPost,
				pullAPIPath,
				fmt.Appendf(nil, `{"title":"%s","head":"%s:%s","base":"%s","body":"%s"}`, title, repo.Owner.Login, head, base, body),
			)

			w.WriteHeader(http.StatusCreated)
			fmt.Fprintf(w, `{"number": %d}\n`, number)
		})

		got, err := g.CreatePull(t.Context(), repo, base, head, title, body)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if got.Number != number {
			t.Fatalf("expected number to be %d but got %d", number, got.Number)
		}

		if got.Repo != repo {
			t.Fatalf("expected repo to be %v but got %v", repo, got.Repo)
		}

		if !got.New {
			t.Fatalf("expected pull to be new, but it is not")
		}
	})
}

func TestGetOrCreatePull(t *testing.T) {
	t.Parallel()

	t.Run("find existing pull", func(t *testing.T) {
		t.Parallel()

		g := setup(t, func(w http.ResponseWriter, r *http.Request) {
			if r.URL.RawQuery != pullAPIQuery {
				t.Fatalf("expected query to be '%s' but got '%s'", pullAPIQuery, r.URL.RawQuery)
			}

			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, `[{"number": %d}]\n`, number)
		})

		got, err := g.GetOrCreatePull(t.Context(), repo, base, head, title, body)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if got.Number != number {
			t.Fatalf("want number %d, but got %d", number, got.Number)
		}

		if got.New {
			t.Fatalf("expected pull to be not new, but it is")
		}
	})

	t.Run("fails to get existing pull", func(t *testing.T) {
		t.Parallel()

		g := setup(t, func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		})

		_, err := g.GetOrCreatePull(t.Context(), repo, base, head, title, body)
		if err == nil {
			t.Fatalf("expected error, got nil")
		}
	})

	t.Run("creates pull if it does not exist", func(t *testing.T) {
		t.Parallel()

		i := 0
		g := setup(t, func(w http.ResponseWriter, r *http.Request) {
			switch i {
			case 0:
				assertReq(t, r, http.MethodGet, pullAPIPath, nil)
				w.WriteHeader(http.StatusOK)
				fmt.Fprintln(w, `[]`)
			case 1:
				assertReq(t, r, http.MethodPost, pullAPIPath, nil)
				w.WriteHeader(http.StatusOK)
				fmt.Fprintf(w, `{"number": %d}\n`, number)
			default:
				w.WriteHeader(http.StatusInternalServerError)
			}

			i++
		})

		got, err := g.GetOrCreatePull(t.Context(), repo, base, head, title, body)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if got.Number != number {
			t.Fatalf("expected number to be %d but got %d", number, got.Number)
		}

		if !got.New {
			t.Fatalf("expected pull to be new, but it is not")
		}
	})
}
