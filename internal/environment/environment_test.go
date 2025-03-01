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
