package config

import (
	_ "embed"
	"regexp"
	"strings"
	"testing"

	"github.com/nobe4/dent.go"

	"github.com/nobe4/action-ln/internal/github"
)

//go:embed all-cases.yaml
var allCases string

func TestConfigParseAll(t *testing.T) {
	t.Parallel()

	c := New()

	err := c.Parse(strings.NewReader(allCases))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	for i, l := range c.Links {
		t.Logf("LINK[%d] %s", i, l.String())
	}

	wants := []string{}

	for _, want := range regexp.
		MustCompile(`(?m)^\s+# want: (.+)$`).
		FindAllStringSubmatch(allCases, -1) {
		wants = append(wants, want[1])
	}

	for i, want := range wants {
		t.Logf("WANT[%d] %s", i, want)
	}

	t.Skip("TODO once the fillmissing is done")

	if ll, lw := len(c.Links), len(wants); ll != lw {
		t.Fatalf("want %d links, but got %d", lw, ll)
	}

	for i, l := range c.Links {
		if l.String() != wants[i] {
			t.Fatalf("want link %d to be %q, but got %q", i, wants[i], l.String())
		}
	}
}

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
			want:   Config{Links: Links{}},
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
				Links: Links{
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
      path: e
      ref: f
`),
			config: &Config{},
			want: Config{
				Links: Links{
					{
						From: github.File{Repo: repo, Path: "c", Ref: "d"},
						To:   github.File{Repo: repo, Path: "e", Ref: "f"},
					},
				},
			},
		},

		// 		{
		// 			name:  "uses defaults",
		// 			input: `links: []`,
		// 			config: &Config{
		// 				Defaults: Defaults{Repo: repo},
		// 			},
		// 			want: Config{
		// 				Defaults: Defaults{Repo: repo},
		// 				Links:    Links{},
		// 			},
		// 		},
		//
		// 		{
		// 			name: "keep defaults",
		// 			input: dent.DedentString(`
		// defaults:
		//   repo: a/b
		// `),
		// 			config: &Config{
		// 				Defaults: Defaults{Repo: repo},
		// 			},
		// 			want: Config{
		// 				Defaults: Defaults{Repo: repo},
		// 				Links:    Links{},
		// 			},
		// 		},
		//
		// 		{
		// 			name: "update defaults",
		// 			input: dent.DedentString(`
		// defaults:
		//   repo: x/y
		// `),
		// 			config: &Config{
		// 				Defaults: Defaults{Repo: repo},
		// 			},
		// 			want: Config{
		// 				Defaults: Defaults{Repo: github.Repo{Owner: github.User{Login: "x"}, Repo: "y"}},
		// 				Links:    Links{},
		// 			},
		// 		},
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
