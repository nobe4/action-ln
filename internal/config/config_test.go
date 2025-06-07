package config

import (
	"embed"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/nobe4/action-ln/internal/github"
)

//go:embed fixtures/*
var fixtures embed.FS

func TestConfigParseAll(t *testing.T) {
	t.Parallel()

	test := func(t *testing.T, path, content string) {
		t.Helper()

		t.Run(path, func(t *testing.T) {
			t.Parallel()

			repo := github.Repo{
				Repo:  "current_repo",
				Owner: github.User{Login: "current_owner"},
			}
			source := github.File{
				Path: ".ln-config.yaml",
				Repo: repo,
			}

			c := New(source, repo)

			err := c.Parse(strings.NewReader(content))
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			for i, l := range c.Links {
				t.Logf("LINK[%d] %s", i, l.String())
			}

			wants := []string{}

			for _, want := range regexp.
				MustCompile(`(?m)^\s+# want: (.+)$`).
				FindAllStringSubmatch(content, -1) {
				wants = append(wants, want[1])
			}

			for i, want := range wants {
				t.Logf("WANT[%d] %s", i, want)
			}

			if ll, lw := len(c.Links), len(wants); ll != lw {
				t.Fatalf("want %d links, but got %d", lw, ll)
			}

			for i, l := range c.Links {
				if l.String() != wants[i] {
					t.Fatalf("want link %d to be %q, but got %q", i, wants[i], l.String())
				}
			}
		})
	}

	fs, err := fixtures.ReadDir("fixtures")
	if err != nil {
		t.Fatalf("failed to list fixtures: %v", err)
	}

	for _, f := range fs {
		path := filepath.Join("fixtures", f.Name())

		content, err := fixtures.ReadFile(path)
		if err != nil {
			t.Fatalf("failed to read fixtures %q: %v", path, err)
		}

		test(t, path, string(content))
	}
}

func TestGetMapKey(t *testing.T) {
	t.Parallel()

	m := map[string]any{
		"a": "a",
		"b": 2,
		"c": []string{"c"},
	}

	if got := getMapKey(m, "a"); got != "a" {
		t.Errorf("want a, but got %v", got)
	}

	if got := getMapKey(m, "b"); got != "" {
		t.Errorf("want \"\", but got %v", got)
	}

	if got := getMapKey(m, "c"); got != "" {
		t.Errorf("want \"\", but got %v", got)
	}
}
