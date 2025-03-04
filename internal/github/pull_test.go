package github

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetPull(t *testing.T) {
	t.Parallel()

	repo := Repo{Owner: User{Login: "owner"}, Repo: "repo"}

	t.Run("finds a PR", func(t *testing.T) {
		t.Parallel()

		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assertReq(t, r, http.MethodGet, "/repos/owner/repo/pulls", nil)

			expectedQuery := "base=base&head=owner%3Ahead&per_page=1&state=open"
			if r.URL.RawQuery != expectedQuery {
				t.Fatalf("expected query to be '%s' but got '%s'", expectedQuery, r.URL.RawQuery)
			}

			w.WriteHeader(http.StatusOK)
			fmt.Fprintln(w, `[{"number": 123}]`)
		}))

		g := New("token", ts.URL)

		got, err := g.GetPull(t.Context(), repo, "base", "head")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if got.Number != 123 {
			t.Fatalf("expected number to be '123' but got %d", got.Number)
		}
	})

	t.Run("finds no PR", func(t *testing.T) {
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
