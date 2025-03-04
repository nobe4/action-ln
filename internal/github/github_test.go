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

// TODO: remove.
func TestGetUser(t *testing.T) {
	t.Parallel()

	t.Run("fails", func(t *testing.T) {
		t.Parallel()

		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusUnauthorized)
		}))

		g := New("token", ts.URL)

		_, err := g.GetUser(t.Context())
		if !errors.Is(err, errRequest) {
			t.Fatalf("expected request error, got %v", err)
		}
	})

	t.Run("succeeds", func(t *testing.T) {
		t.Parallel()

		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			fmt.Fprintln(w, `{"login": "user"}`)
		}))

		g := New("token", ts.URL)

		u, err := g.GetUser(t.Context())
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if u.Login != "user" {
			t.Fatal("expected user to parse correctly")
		}
	})
}

func TestReq(t *testing.T) {
	t.Parallel()

	t.Run("fails to authenticate", func(t *testing.T) {
		t.Parallel()

		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusUnauthorized)
		}))

		g := New("token", ts.URL)

		status, err := g.req(t.Context(), http.MethodGet, PathUser, nil, nil)
		if !errors.Is(err, errRequest) {
			t.Fatalf("expected request error, got %v", err)
		}

		if status != http.StatusUnauthorized {
			t.Fatalf("expected %d, got %d", http.StatusUnauthorized, status)
		}
	})

	t.Run("fails with 500", func(t *testing.T) {
		t.Parallel()

		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))

		g := New("token", ts.URL)

		status, err := g.req(t.Context(), http.MethodGet, PathUser, nil, nil)
		if !errors.Is(err, errRequest) {
			t.Fatalf("expected request error, got %v", err)
		}

		if status != http.StatusInternalServerError {
			t.Fatalf("expected %d, got %d", http.StatusInternalServerError, status)
		}
	})

	t.Run("fails to decode JSON", func(t *testing.T) {
		t.Parallel()

		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			fmt.Fprintln(w, `<invalid json>`)
		}))

		g := New("token", ts.URL)
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

		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assertReq(t, r, http.MethodGet, PathUser, nil)

			if auth := r.Header.Get("Authorization"); auth != "Bearer token" {
				t.Fatal("invalid token", auth)
			}

			w.WriteHeader(http.StatusOK)
			fmt.Fprintln(w, `{"data":"123"}`)
		}))

		g := New("token", ts.URL)

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

		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			fmt.Fprintln(w, `{"success": true}`)
		}))

		g := New("token", ts.URL)
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
