package config

import (
	"fmt"
	"testing"

	"github.com/nobe4/action-ln/internal/github"
)

//nolint:maintidx // This is just a big list of tests.
func TestParseFile(t *testing.T) {
	t.Parallel()

	const complexPath = "a/b-c/d_f/f.txt"

	tests := []struct {
		input any
		want  []github.File
	}{
		// nil
		{},

		// Slice
		{
			input: []any{nil},
			want:  []github.File{},
		},

		{
			input: []any{nil, nil, nil},
			want:  []github.File{},
		},

		{
			input: []any{
				map[string]any{"path": "path"},
				"path2",
			},
			want: []github.File{
				{Path: "path"},
				{Path: "path2"},
			},
		},

		{
			input: []any{
				map[string]any{"path": "path"},
				"path2",
			},
			want: []github.File{
				{Path: "path"},
				{Path: "path2"},
			},
		},

		// Map
		{
			input: map[string]any{"path": "path"},
			want:  []github.File{{Path: "path"}},
		},

		{
			input: map[string]any{"path": "path"},
			want:  []github.File{{Path: "path"}},
		},

		{
			input: map[string]any{"repo": "repo", "path": "path"},
			want: []github.File{
				{
					Repo: github.Repo{
						Repo: "repo",
					},
					Path: "path",
				},
			},
		},

		{
			input: map[string]any{"repo": "repo2", "path": "path"},
			want: []github.File{
				{
					Repo: github.Repo{Repo: "repo2"},
					Path: "path",
				},
			},
		},

		{
			input: map[string]any{"repo": "repo", "owner": "owner", "path": "path", "ref": "ref"},
			want: []github.File{
				{
					Repo: github.Repo{
						Owner: github.User{Login: "owner"},
						Repo:  "repo",
					},
					Path: "path",
					Ref:  "ref",
				},
			},
		},

		{
			input: map[string]any{"repo": "repo/owner", "path": "path"},
			want: []github.File{
				{
					Repo: github.Repo{
						Owner: github.User{Login: "repo"},
						Repo:  "owner",
					},
					Path: "path",
				},
			},
		},

		// String
		{
			input: "https://github.com/owner/repo/blob/ref/path",
			want: []github.File{
				{
					Repo: github.Repo{
						Owner: github.User{Login: "owner"},
						Repo:  "repo",
					},
					Path: "path",
					Ref:  "ref",
				},
			},
		},

		{
			input: "https://github.com/owner/repo/blob/ref/path",
			want: []github.File{
				{
					Repo: github.Repo{
						Owner: github.User{Login: "owner"},
						Repo:  "repo",
					},
					Path: "path",
					Ref:  "ref",
				},
			},
		},

		{
			input: "https://github.com/owner/repo/blob/ref/" + complexPath,
			want: []github.File{
				{
					Repo: github.Repo{
						Owner: github.User{Login: "owner"},
						Repo:  "repo",
					},
					Path: complexPath,
					Ref:  "ref",
				},
			},
		},

		{
			input: "owner/repo/blob/ref/path",
			want: []github.File{
				{
					Repo: github.Repo{
						Owner: github.User{Login: "owner"},
						Repo:  "repo",
					},
					Path: "path",
					Ref:  "ref",
				},
			},
		},

		{
			input: "owner/repo/blob/ref/path",
			want: []github.File{
				{
					Repo: github.Repo{
						Owner: github.User{Login: "owner"},
						Repo:  "repo",
					},
					Path: "path",
					Ref:  "ref",
				},
			},
		},

		{
			input: "owner/repo/blob/ref/" + complexPath,
			want: []github.File{
				{
					Repo: github.Repo{
						Owner: github.User{Login: "owner"},
						Repo:  "repo",
					},
					Path: complexPath,
					Ref:  "ref",
				},
			},
		},

		{
			input: "owner/repo:path@ref",
			want: []github.File{
				{
					Repo: github.Repo{
						Owner: github.User{Login: "owner"},
						Repo:  "repo",
					},
					Path: "path",
					Ref:  "ref",
				},
			},
		},

		{
			input: "owner/repo:path@ref",
			want: []github.File{
				{
					Repo: github.Repo{
						Owner: github.User{Login: "owner"},
						Repo:  "repo",
					},
					Path: "path",
					Ref:  "ref",
				},
			},
		},

		{
			input: "owner/repo:" + complexPath + "@ref",
			want: []github.File{
				{
					Repo: github.Repo{
						Owner: github.User{Login: "owner"},
						Repo:  "repo",
					},
					Path: complexPath,
					Ref:  "ref",
				},
			},
		},

		{
			input: "owner/:path@ref",
			want: []github.File{
				{
					Repo: github.Repo{
						Owner: github.User{Login: "owner"},
					},
					Path: "path",
					Ref:  "ref",
				},
			},
		},

		{
			input: "/repo:path@ref",
			want: []github.File{
				{
					Repo: github.Repo{
						Repo: "repo",
					},
					Path: "path",
					Ref:  "ref",
				},
			},
		},

		{
			input: "owner/:",
			want: []github.File{
				{
					Repo: github.Repo{
						Owner: github.User{Login: "owner"},
					},
				},
			},
		},

		{
			input: "/repo:",
			want: []github.File{
				{
					Repo: github.Repo{
						Repo: "repo",
					},
				},
			},
		},

		{
			input: "path@ref",
			want:  []github.File{{Path: "path", Ref: "ref"}},
		},

		{
			input: "path@ref",
			want:  []github.File{{Path: "path", Ref: "ref"}},
		},

		{
			input: complexPath + "@ref",
			want:  []github.File{{Path: complexPath, Ref: "ref"}},
		},

		{
			input: "path",
			want:  []github.File{{Path: "path"}},
		},

		{
			input: "path",
			want:  []github.File{{Path: "path"}},
		},

		{
			input: complexPath,
			want:  []github.File{{Path: complexPath}},
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("%v", test.input), func(t *testing.T) {
			t.Parallel()

			c := New()

			got, err := c.parseFile(test.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if lg, lw := len(got), len(test.want); lg != lw {
				t.Fatalf("want %d files, but got %d", lw, lg)
			}

			for i, f := range got {
				if !f.Equal(test.want[i]) {
					t.Errorf("file %d: want %+v, but got %+v", i, test.want, got)
				}
			}
		})
	}
}
