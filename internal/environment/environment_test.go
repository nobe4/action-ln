package environment

import (
	"errors"
	"testing"
)

func TestParseNoop(t *testing.T) {
	tests := []struct {
		noop string
		want bool
	}{
		{noop: "", want: false},
		{noop: "no", want: false},
		{noop: "false", want: false},
		{noop: "0", want: false},
		{noop: "yes", want: false},
		{noop: "1", want: false},
		{noop: "true", want: true},
	}

	for _, test := range tests {
		t.Run(test.noop, func(t *testing.T) {
			t.Setenv("INPUT_NOOP", test.noop)

			got := parseNoop()
			if got != test.want {
				t.Fatalf("want %v but got %v", test.want, got)
			}
		})
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

		endpoint := parseEndpoint()
		if want != endpoint {
			t.Fatalf("want %v but got %v", want, endpoint)
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

	t.Run("gets the set endpoint", func(t *testing.T) {
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
