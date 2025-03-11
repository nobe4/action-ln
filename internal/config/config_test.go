package config

import (
	"strings"
	"testing"

	"github.com/nobe4/dent.go"

	"github.com/nobe4/action-ln/internal/github"
)

func TestConfigParse(t *testing.T) {
	t.Parallel()

	repo := github.Repo{Owner: github.User{Login: "a"}, Repo: "b"}

	tests := []struct {
		name   string
		input  string
		config *Config
		want   Config
	}{
		{
			name:   "empty",
			input:  `links: []`,
			config: &Config{},
			want:   Config{Links: []Link{}},
		},

		{
			name: "path-only link",
			input: dent.DedentString(`
links:
  - from: a
    to: b
`),
			config: &Config{},
			want: Config{
				Links: []Link{
					{
						From: github.File{Path: "a"},
						To:   github.File{Path: "b"},
					},
				},
			},
		},

		{
			name: "complex link",
			input: dent.DedentString(`
links:
  - from: a/b:c@d
    to:
      owner: a
      repo: b
      path: c
      ref: d
`),
			config: &Config{},
			want: Config{
				Links: []Link{
					{
						From: github.File{Repo: repo, Path: "c", Ref: "d"},
						To:   github.File{Repo: repo, Path: "c", Ref: "d"},
					},
				},
			},
		},

		{
			name:  "uses defaults",
			input: `links: []`,
			config: &Config{
				Defaults: Defaults{Repo: repo},
			},
			want: Config{
				Defaults: Defaults{Repo: repo},
				Links:    []Link{},
			},
		},

		{
			name: "keep defaults",
			input: dent.DedentString(`
defaults:
  repo: a/b
`),
			config: &Config{
				Defaults: Defaults{Repo: repo},
			},
			want: Config{
				Defaults: Defaults{Repo: repo},
				Links:    []Link{},
			},
		},

		{
			name: "update defaults",
			input: dent.DedentString(`
defaults:
  repo: x/y
`),
			config: &Config{
				Defaults: Defaults{Repo: repo},
			},
			want: Config{
				Defaults: Defaults{Repo: github.Repo{Owner: github.User{Login: "x"}, Repo: "y"}},
				Links:    []Link{},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			err := test.config.Parse(strings.NewReader(test.input))
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}

			if test.config.String() != test.want.String() {
				t.Errorf("want %#v, got %#v", test.want, test.config)
			}
		})
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
