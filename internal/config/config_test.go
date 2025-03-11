package config

import (
	"strings"
	"testing"

	"github.com/nobe4/dent.go"

	"github.com/nobe4/action-ln/internal/environment"
	"github.com/nobe4/action-ln/internal/github"
)

func TestConfigParse(t *testing.T) {
	t.Parallel()

	tests := []struct {
		input string
		env   environment.Environment
		want  Config
	}{
		{
			input: `links: []`,
			want: Config{
				Links: []Link{},
			},
		},

		{
			input: dent.DedentString(`
links:
  - from: a
    to: b
`),
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
			input: dent.DedentString(`
links:
  - from: a/b:c@d
    to:
      owner: e
      repo: f
      path: g
      ref: h
`),
			want: Config{
				Links: []Link{
					{
						From: github.File{
							Repo: github.Repo{Owner: github.User{Login: "a"}, Repo: "b"},
							Path: "c",
							Ref:  "d",
						},
						To: github.File{
							Repo: github.Repo{Owner: github.User{Login: "e"}, Repo: "f"},
							Path: "g",
							Ref:  "h",
						},
					},
				},
			},
		},

		{
			input: `links: []`,
			env: environment.Environment{
				Repo: github.Repo{Owner: github.User{Login: "a"}, Repo: "b"},
			},
			want: Config{
				Defaults: Defaults{
					Repo: github.Repo{Owner: github.User{Login: "a"}, Repo: "b"},
				},
				Links: []Link{},
			},
		},

		{
			input: dent.DedentString(`
defaults:
  repo: a/b
`),
			env: environment.Environment{
				Repo: github.Repo{Owner: github.User{Login: "x"}, Repo: "y"},
			},
			want: Config{
				Defaults: Defaults{
					Repo: github.Repo{Owner: github.User{Login: "a"}, Repo: "b"},
				},
				Links: []Link{},
			},
		},

		{
			input: dent.DedentString(`
defaults:
  repo: a/b
`),
			want: Config{
				Defaults: Defaults{
					Repo: github.Repo{Owner: github.User{Login: "a"}, Repo: "b"},
				},
				Links: []Link{},
			},
		},

		{
			input: dent.DedentString(`
defaults:
  owner: a
  repo: b
`),
			want: Config{
				Defaults: Defaults{
					Repo: github.Repo{Owner: github.User{Login: "a"}, Repo: "b"},
				},
				Links: []Link{},
			},
		},
	}

	for _, test := range tests {
		t.Run("", func(t *testing.T) {
			t.Parallel()

			c := Config{}

			err := c.Parse(strings.NewReader(test.input), test.env)
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}

			if c.String() != test.want.String() {
				t.Errorf("want %v, got %v", test.want, c)
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
