package environment

import (
	"errors"
	"testing"
)

func TestParseToken(t *testing.T) {
	const want = "token"

	t.Run("gets INPUT_TOKEN", func(t *testing.T) {
		t.Setenv("INPUT_TOKEN", want)

		got, err := parseToken()
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if got != want {
			t.Fatalf("want %v, got %v", want, got)
		}
	})

	t.Run("gets GITHUB_TOKEN", func(t *testing.T) {
		t.Setenv("GITHUB_TOKEN", want)

		got, err := parseToken()
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if got != want {
			t.Fatalf("want %v, got %v", want, got)
		}
	})

	t.Run("gets nothing", func(t *testing.T) {
		// Need to force an empty value to not conflict with GitHub Action's Env
		t.Setenv("INPUT_TOKEN", "")

		_, err := parseToken()
		if !errors.Is(err, errNoToken) {
			t.Fatalf("want %v but got error: %v", errNoToken, err)
		}
	})
}

func TestParseRepo(t *testing.T) {
	t.Run("gets nothing", func(t *testing.T) {
		// Need to force an empty value to not conflict with GitHub Action's Env
		t.Setenv("GITHUB_REPOSITORY", "")

		_, err := parseRepo()
		if !errors.Is(err, errNoRepo) {
			t.Fatalf("want %v but got error: %v", errNoRepo, err)
		}
	})

	t.Run("gets nothing", func(t *testing.T) {
		t.Setenv("GITHUB_REPOSITORY", "owner+repo+is+invalid")

		_, err := parseRepo()
		if !errors.Is(err, errInvalidRepo) {
			t.Fatalf("want %v but got error: %v", errInvalidRepo, err)
		}
	})

	t.Run("gets the parsed Repo", func(t *testing.T) {
		t.Setenv("GITHUB_REPOSITORY", "owner/repo")

		got, err := parseRepo()
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if got.Owner.Login != "owner" || got.Repo != "repo" {
			t.Fatalf("want %v but got %+v", "owner/repo", got)
		}
	})
}

func TestParseEndpoint(t *testing.T) {
	t.Run("gets the default", func(t *testing.T) {
		// Need to force an empty value to not conflict with GitHub Action's Env
		t.Setenv("GITHUB_REPOSITORY", "")

		env := parseEndpoint()
		if defaultEndpoint != env {
			t.Fatalf("want %v but got %v", defaultEndpoint, env)
		}
	})

	t.Run("gets the set endpoint", func(t *testing.T) {
		want := "https://example.com"
		t.Setenv("GITHUB_API_URL", want)

		endpoint := parseEndpoint()
		if want != endpoint {
			t.Fatalf("want %v but got %v", want, endpoint)
		}
	})
}
