package github

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetUser(t *testing.T) {
	t.Parallel()

	t.Run("fails", func(t *testing.T) {
		t.Parallel()

		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/user" {
				t.Error("expected request to /user")
			}

			w.WriteHeader(http.StatusUnauthorized)
		}))

		g := New("token", ts.URL)

		err := g.req(t.Context(), "GET", "/user", nil, nil)
		if !errors.Is(err, errRequest) {
			t.Fatalf("expected request error, got %v", err)
		}
	})

	t.Run("succeeds", func(t *testing.T) {
		t.Parallel()

		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/user" {
				t.Error("expected request to /user")
			}

			w.WriteHeader(http.StatusOK)
			fmt.Fprintln(w, `{"login": "user"}`)
		}))

		user := User{}
		g := New("token", ts.URL)

		err := g.req(t.Context(), "GET", "/user", nil, &user)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if user.Login != "user" {
			t.Fatal("expected user to parse correctly")
		}
	})
}

func TestReq(t *testing.T) {
	t.Parallel()

	t.Run("fails to authenticate", func(t *testing.T) {
		t.Parallel()

		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/user" {
				t.Error("expected request to /user")
			}

			if r.Header.Get("Authorization") != "Bearer token" {
				t.Error("expected Authorization header to be set")
			}

			w.WriteHeader(http.StatusUnauthorized)
		}))

		g := New("token", ts.URL)

		err := g.req(t.Context(), "GET", "/user", nil, nil)
		if !errors.Is(err, errRequest) {
			t.Fatalf("expected request error, got %v", err)
		}
	})

	t.Run("fails with 500", func(t *testing.T) {
		t.Parallel()

		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))

		g := New("token", ts.URL)

		err := g.req(t.Context(), "GET", "/user", nil, nil)
		if !errors.Is(err, errRequest) {
			t.Fatalf("expected request error, got %v", err)
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

		err := g.req(t.Context(), "GET", "/user", nil, &data)

		var jsonErr *json.SyntaxError
		if !errors.As(err, &jsonErr) {
			t.Fatalf("expected json syntax error, got %v", err)
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

		err := g.req(t.Context(), "GET", "/user", nil, &data)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if !data.Success {
			t.Fatal("expected success")
		}
	})
}
