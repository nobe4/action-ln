package environment

import (
	"errors"
	"testing"
)

//nolint:tparallel // t.Setenv is not thread-safe.
func TestParseToken(t *testing.T) {
	const want = "token"

	t.Run("gets INPUT_TOKEN", func(t *testing.T) {
		t.Setenv("INPUT_TOKEN", want)

		got, err := parseToken()
		if err != nil {
			t.Fatalf("got error: %v", err)
		}

		if got != want {
			t.Fatalf("want %v, got %v", want, got)
		}
	})

	t.Run("gets GITHUB_TOKEN", func(t *testing.T) {
		t.Setenv("GITHUB_TOKEN", want)

		got, err := parseToken()
		if err != nil {
			t.Fatalf("got error: %v", err)
		}

		if got != want {
			t.Fatalf("want %v, got %v", want, got)
		}
	})

	t.Run("gets nothing", func(t *testing.T) {
		t.Parallel()

		_, err := parseToken()
		if !errors.Is(err, errNoToken) {
			t.Fatalf("want %v but got error: %v", errNoToken, err)
		}
	})
}

//nolint:tparallel // t.Setenv is not thread-safe.
func TestParseRepo(t *testing.T) {
	t.Run("gets nothing", func(t *testing.T) {
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
			t.Fatalf("got error: %v", err)
		}

		if got.Owner.Login != "owner" || got.Repo != "repo" {
			t.Fatalf("want %v but got %+v", "owner/repo", got)
		}
	})
}
