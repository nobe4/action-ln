package github

import (
	"fmt"
	"net/http"
	"testing"
)

const (
	//nolint:gosec // This is not a secret.
	appTokenAPIPath = "/app/installations/123/access_tokens"
	jwt             = "jwt"
	install         = "123"
)

func TestGetAppToken(t *testing.T) {
	t.Parallel()

	g := setup(t, func(w http.ResponseWriter, r *http.Request) {
		assertReq(t, r, http.MethodPost, appTokenAPIPath, nil)

		if auth := r.Header.Get("Authorization"); auth != "Bearer "+jwt {
			t.Fatal("invalid jwt", auth)
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"token": "%s"}`, token)
	})

	got, err := g.GetAppToken(t.Context(), install, jwt)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if got != token {
		t.Fatalf("expected token to be '%s' but got '%s'", token, got)
	}
}
