package github

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

const (
	token = "token"
)

//nolint:gochecknoglobals // This is used across GitHub tests.
var repo = Repo{Owner: User{Login: "owner"}, Repo: "repo"}

func assertReq(t *testing.T, r *http.Request, method, path string, body []byte) {
	t.Helper()

	if r.URL.Path != path {
		t.Fatalf("want path '%s', got %s", path, r.URL.Path)
	}

	if r.Method != method {
		t.Fatalf("want method '%s', got %s", method, r.Method)
	}

	if body != nil {
		gotBody, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatal("failed to read body", err)
		}

		if !bytes.Equal(gotBody, body) {
			t.Fatalf("want body '%s', got '%s'", string(body), string(gotBody))
		}
	}
}

func setup(t *testing.T, f func(w http.ResponseWriter, r *http.Request)) *GitHub {
	t.Helper()

	ts := httptest.NewServer(http.HandlerFunc(f))

	// NOTE: using http.DefaultClient here is expected, as we mock the server
	// with ts.
	g := New(http.DefaultClient, ts.URL)
	g.Token = token

	return g
}

func TestReq(t *testing.T) {
	t.Parallel()

	t.Run("fails to authenticate", func(t *testing.T) {
		t.Parallel()

		g := setup(t, func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusUnauthorized)
		})

		status, err := g.req(t.Context(), http.MethodGet, PathUser, nil, nil)
		if !errors.Is(err, ErrRequestFailed) {
			t.Fatalf("expected request error, got %v", err)
		}

		if status != http.StatusUnauthorized {
			t.Fatalf("expected %d, got %d", http.StatusUnauthorized, status)
		}
	})

	t.Run("fails with 500", func(t *testing.T) {
		t.Parallel()

		g := setup(t, func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		})

		status, err := g.req(t.Context(), http.MethodGet, PathUser, nil, nil)
		if !errors.Is(err, ErrRequestFailed) {
			t.Fatalf("expected request error, got %v", err)
		}

		if status != http.StatusInternalServerError {
			t.Fatalf("expected %d, got %d", http.StatusInternalServerError, status)
		}
	})

	t.Run("fails to decode JSON", func(t *testing.T) {
		t.Parallel()

		g := setup(t, func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusOK)
			fmt.Fprintln(w, `<invalid json>`)
		})

		data := ""

		status, err := g.req(t.Context(), http.MethodGet, PathUser, nil, &data)

		var jsonErr *json.SyntaxError
		if !errors.As(err, &jsonErr) {
			t.Fatalf("expected json syntax error, got %v", err)
		}

		if status != http.StatusInternalServerError {
			t.Fatalf("expected %d, got %d", http.StatusInternalServerError, status)
		}
	})

	t.Run("decodes nothing", func(t *testing.T) {
		t.Parallel()

		g := setup(t, func(w http.ResponseWriter, r *http.Request) {
			assertReq(t, r, http.MethodGet, PathUser, nil)

			if auth := r.Header.Get("Authorization"); auth != "Bearer token" {
				t.Fatal("invalid token", auth)
			}

			w.WriteHeader(http.StatusOK)
			fmt.Fprintln(w, `{"data":"123"}`)
		})

		status, err := g.req(t.Context(), http.MethodGet, PathUser, nil, nil)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if status != http.StatusOK {
			t.Fatalf("expected %d, got %d", http.StatusOK, status)
		}
	})

	t.Run("decodes JSON response correctly", func(t *testing.T) {
		t.Parallel()

		g := setup(t, func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusOK)
			fmt.Fprintln(w, `{"success": true}`)
		})

		data := struct{ Success bool }{}

		status, err := g.req(t.Context(), http.MethodGet, PathUser, nil, &data)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if !data.Success {
			t.Fatal("expected success")
		}

		if status != http.StatusOK {
			t.Fatalf("expected %d, got %d", http.StatusOK, status)
		}
	})
}
