package environment

import (
	"errors"
	"testing"

	"github.com/nobe4/gh-ln/pkg/environment"
)

func TestParseNoop(t *testing.T) {
	t.Setenv("INPUT_NOOP", "")

	if parseNoop() {
		t.Fatal("want false but got true")
	}

	t.Setenv("INPUT_NOOP", "1")

	if !parseNoop() {
		t.Fatal("want true but got false")
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
		if !errors.Is(err, environment.ErrNoToken) {
			t.Fatalf("want %v but got error: %v", environment.ErrNoToken, err)
		}
	})
}

func TestParseEndpoint(t *testing.T) {
	t.Run("gets the default", func(t *testing.T) {
		// Need to force an empty value to not conflict with GitHub Action's Env
		t.Setenv("GITHUB_REPOSITORY", "")

		got := parseEndpoint()
		if environment.DefaultEndpoint != got {
			t.Fatalf("want %v but got %v", environment.DefaultEndpoint, got)
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
		if environment.DefaultServer != got {
			t.Fatalf("want %v but got %v", environment.DefaultServer, got)
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
		if environment.DefaultRunID != got {
			t.Fatalf("want %v but got %v", environment.DefaultRunID, got)
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
		if environment.DefaultConfig != got {
			t.Fatalf("want %v but got %v", environment.DefaultConfig, got)
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
		t.Fatal("want false but got true")
	}

	t.Setenv("GITHUB_RUN_ID", "1234")

	if !parseOnAction() {
		t.Fatal("want true but got false")
	}
}

func TestParseDebug(t *testing.T) {
	t.Setenv("RUNNER_DEBUG", "")

	if parseDebug() {
		t.Fatal("want false but got true")
	}

	t.Setenv("RUNNER_DEBUG", "1")

	if !parseDebug() {
		t.Fatal("want true but got false")
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
