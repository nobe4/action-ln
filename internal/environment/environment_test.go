package environment

import (
	"errors"
	"testing"
)

func TestParseNoop(t *testing.T) {
	t.Setenv("INPUT_NOOP", "")

	if parseNoop() {
		t.Fatalf("want false but got true")
	}

	t.Setenv("INPUT_NOOP", "1")

	if !parseNoop() {
		t.Fatalf("want true but got false")
	}
}

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
		if !errors.Is(err, ErrNoToken) {
			t.Fatalf("want %v but got error: %v", ErrNoToken, err)
		}
	})
}

func TestParseRepo(t *testing.T) {
	t.Run("gets nothing", func(t *testing.T) {
		// Need to force an empty value to not conflict with GitHub Action's Env
		t.Setenv("GITHUB_REPOSITORY", "")

		_, err := parseRepo()
		if !errors.Is(err, ErrNoRepo) {
			t.Fatalf("want %v but got error: %v", ErrNoRepo, err)
		}
	})

	t.Run("gets nothing", func(t *testing.T) {
		t.Setenv("GITHUB_REPOSITORY", "owner+repo+is+invalid")

		_, err := parseRepo()
		if !errors.Is(err, ErrInvalidRepo) {
			t.Fatalf("want %v but got error: %v", ErrInvalidRepo, err)
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

		got := parseEndpoint()
		if defaultEndpoint != got {
			t.Fatalf("want %v but got %v", defaultEndpoint, got)
		}
	})

	t.Run("gets the set endpoint", func(t *testing.T) {
		want := "https://example.com"
		t.Setenv("GITHUB_API_URL", want)

		got := parseEndpoint()
		if want != got {
			t.Fatalf("want %v but got %v", want, got)
		}
	})
}

func TestParseServer(t *testing.T) {
	t.Run("gets the default", func(t *testing.T) {
		// Need to force an empty value to not conflict with GitHub Action's Env
		t.Setenv("GITHUB_SERVER_URL", "")

		got := parseServer()
		if defaultServer != got {
			t.Fatalf("want %v but got %v", defaultServer, got)
		}
	})

	t.Run("gets the set server", func(t *testing.T) {
		want := "https://example.com"
		t.Setenv("GITHUB_SERVER_URL", want)

		got := parseServer()
		if want != got {
			t.Fatalf("want %v but got %v", want, got)
		}
	})
}

func TestParseRunID(t *testing.T) {
	t.Run("gets the default", func(t *testing.T) {
		// Need to force an empty value to not conflict with GitHub Action's Env
		t.Setenv("GITHUB_RUN_ID", "")

		got := parseRunID()
		if defaultRunID != got {
			t.Fatalf("want %v but got %v", defaultRunID, got)
		}
	})

	t.Run("gets the set runID", func(t *testing.T) {
		want := "runID"
		t.Setenv("GITHUB_RUN_ID", want)

		got := parseRunID()
		if want != got {
			t.Fatalf("want %v but got %v", want, got)
		}
	})
}

func TestParseConfig(t *testing.T) {
	t.Run("gets the default", func(t *testing.T) {
		// Need to force an empty value to not conflict with GitHub Action's Env
		t.Setenv("INPUT_CONFIG", "")

		got := parseConfig()
		if defaultConfig != got {
			t.Fatalf("want %v but got %v", defaultConfig, got)
		}
	})

	t.Run("gets the set config", func(t *testing.T) {
		want := "path/to/config"
		t.Setenv("INPUT_CONFIG", want)

		got := parseConfig()
		if want != got {
			t.Fatalf("want %v but got %v", want, got)
		}
	})
}

func TestParseApp(t *testing.T) {
	want := "value"
	t.Setenv("INPUT_APP_ID", want)
	t.Setenv("INPUT_APP_PRIVATE_KEY", want)
	t.Setenv("INPUT_APP_INSTALL_ID", want)

	got := parseApp()

	if want != got.ID {
		t.Fatalf("want %v but got %v", want, got)
	}

	if want != got.PrivateKey {
		t.Fatalf("want %v but got %v", want, got)
	}

	if want != got.InstallID {
		t.Fatalf("want %v but got %v", want, got)
	}
}

func TestParseOnAction(t *testing.T) {
	t.Setenv("GITHUB_RUN_ID", "")

	if parseOnAction() {
		t.Fatalf("want false but got true")
	}

	t.Setenv("GITHUB_RUN_ID", "1234")

	if !parseOnAction() {
		t.Fatalf("want true but got false")
	}
}

func TestParseDebug(t *testing.T) {
	t.Setenv("RUNNER_DEBUG", "")

	if parseDebug() {
		t.Fatalf("want false but got true")
	}

	t.Setenv("RUNNER_DEBUG", "1")

	if !parseDebug() {
		t.Fatalf("want true but got false")
	}
}

func TestTruthy(t *testing.T) {
	t.Parallel()

	tests := []struct {
		s    string
		want bool
	}{
		{s: "", want: false},
		{s: "1", want: true},
		{s: "123", want: false},
		{s: "true", want: true},
		{s: "True", want: true},
		{s: "TRUE", want: true},
		{s: "yes", want: true},
		{s: "Yes", want: true},
		{s: "YES", want: true},
	}

	for _, test := range tests {
		t.Run("", func(t *testing.T) {
			t.Parallel()

			if truthy(test.s) != test.want {
				t.Fatalf("want %v but got %v", test.want, truthy(test.s))
			}
		})
	}
}

func TestMissingOrRedacted(t *testing.T) {
	t.Parallel()

	if missingOrRedacted("token") != redacted {
		t.Fatalf("want redacted but got something else")
	}

	if missingOrRedacted("") != missing {
		t.Fatalf("want missingd but got something else")
	}
}
