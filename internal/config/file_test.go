package config

import (
	"fmt"
	"testing"

	"github.com/nobe4/action-ln/internal/github"
)

func TestParseFileMap(t *testing.T) {
	t.Parallel()

	repo := github.Repo{Owner: github.User{Login: "owner"}, Repo: "repo"}

	tests := []struct {
		defaults Defaults
		input    map[string]any
		want     github.File
	}{
		{},

		{
			input: map[string]any{"path": "path"},
			want:  github.File{Path: "path"},
		},

		{
			defaults: Defaults{Repo: repo},
			input:    map[string]any{"path": "path"},
			want:     github.File{Path: "path", Repo: repo},
		},

		{
			input: map[string]any{"repo": "repo", "path": "path"},
			want: github.File{
				Repo: github.Repo{
					Repo: "repo",
				},
				Path: "path",
			},
		},

		// TODO: might want to inherite the owner from the defaults
		{
			defaults: Defaults{Repo: repo},
			input:    map[string]any{"repo": "repo2", "path": "path"},
			want: github.File{
				Repo: github.Repo{Repo: "repo2"},
				Path: "path",
			},
		},

		{
			input: map[string]any{"repo": "repo", "owner": "owner", "path": "path", "ref": "ref"},
			want: github.File{
				Repo: github.Repo{
					Owner: github.User{Login: "owner"},
					Repo:  "repo",
				},
				Path: "path",
				Ref:  "ref",
			},
		},

		{
			input: map[string]any{"repo": "repo/owner", "path": "path"},
			want: github.File{
				Repo: github.Repo{
					Owner: github.User{Login: "repo"},
					Repo:  "owner",
				},
				Path: "path",
			},
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("%v", test.input), func(t *testing.T) {
			t.Parallel()

			c := New()
			c.Defaults = test.defaults

			got, err := c.parseFileMap(test.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if !test.want.Equal(got) {
				t.Errorf("want %+v, but got %+v", test.want, got)
			}
		})
	}
}

func TestParseFileString(t *testing.T) {
	t.Parallel()

	const complexPath = "a/b-c/d_f/f.txt"

	repo := github.Repo{Owner: github.User{Login: "owner test"}, Repo: "repo test"}

	tests := []struct {
		defaults Defaults
		input    string
		want     github.File
	}{
		{
			input: "https://github.com/owner/repo/blob/ref/path",
			want: github.File{
				Repo: github.Repo{
					Owner: github.User{Login: "owner"},
					Repo:  "repo",
				},
				Path: "path",
				Ref:  "ref",
			},
		},

		{
			defaults: Defaults{Repo: repo},
			input:    "https://github.com/owner/repo/blob/ref/path",
			want: github.File{
				Repo: github.Repo{
					Owner: github.User{Login: "owner"},
					Repo:  "repo",
				},
				Path: "path",
				Ref:  "ref",
			},
		},

		{
			input: "https://github.com/owner/repo/blob/ref/" + complexPath,
			want: github.File{
				Repo: github.Repo{
					Owner: github.User{Login: "owner"},
					Repo:  "repo",
				},
				Path: complexPath,
				Ref:  "ref",
			},
		},

		{
			input: "owner/repo/blob/ref/path",
			want: github.File{
				Repo: github.Repo{
					Owner: github.User{Login: "owner"},
					Repo:  "repo",
				},
				Path: "path",
				Ref:  "ref",
			},
		},

		{
			defaults: Defaults{Repo: repo},
			input:    "owner/repo/blob/ref/path",
			want: github.File{
				Repo: github.Repo{
					Owner: github.User{Login: "owner"},
					Repo:  "repo",
				},
				Path: "path",
				Ref:  "ref",
			},
		},

		{
			input: "owner/repo/blob/ref/" + complexPath,
			want: github.File{
				Repo: github.Repo{
					Owner: github.User{Login: "owner"},
					Repo:  "repo",
				},
				Path: complexPath,
				Ref:  "ref",
			},
		},

		{
			input: "owner/repo:path@ref",
			want: github.File{
				Repo: github.Repo{
					Owner: github.User{Login: "owner"},
					Repo:  "repo",
				},
				Path: "path",
				Ref:  "ref",
			},
		},

		{
			defaults: Defaults{Repo: repo},
			input:    "owner/repo:path@ref",
			want: github.File{
				Repo: github.Repo{
					Owner: github.User{Login: "owner"},
					Repo:  "repo",
				},
				Path: "path",
				Ref:  "ref",
			},
		},

		{
			input: "owner/repo:" + complexPath + "@ref",
			want: github.File{
				Repo: github.Repo{
					Owner: github.User{Login: "owner"},
					Repo:  "repo",
				},
				Path: complexPath,
				Ref:  "ref",
			},
		},

		{
			input: "path@ref",
			want:  github.File{Path: "path", Ref: "ref"},
		},

		{
			defaults: Defaults{Repo: repo},
			input:    "path@ref",
			want:     github.File{Path: "path", Ref: "ref", Repo: repo},
		},

		{
			input: complexPath + "@ref",
			want:  github.File{Path: complexPath, Ref: "ref"},
		},

		{
			input: "path",
			want:  github.File{Path: "path"},
		},

		{
			defaults: Defaults{Repo: repo},
			input:    "path",
			want:     github.File{Path: "path", Repo: repo},
		},

		{
			input: complexPath,
			want:  github.File{Path: complexPath},
		},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			t.Parallel()

			c := New()
			c.Defaults = test.defaults

			got, err := c.parseFileString(test.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if !test.want.Equal(got) {
				t.Errorf("want %+v, but got %+v", test.want, got)
			}
		})
	}
}
