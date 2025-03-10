package config

import (
	"fmt"
	"testing"

	"github.com/nobe4/action-ln/internal/github"
)

func TestParseFileMap(t *testing.T) {
	t.Parallel()

	tests := []struct {
		input map[string]any
		want  github.File
	}{
		{},
		{
			input: map[string]any{"path": "z"},
			want:  github.File{Path: "z"},
		},
		{
			input: map[string]any{"repo": "x", "path": "z"},
			want: github.File{
				Repo: github.Repo{
					Repo: "x",
				},
				Path: "z",
			},
		},
		{
			input: map[string]any{"repo": "x", "owner": "y", "path": "z", "ref": "r"},
			want: github.File{
				Repo: github.Repo{
					Owner: github.User{Login: "y"},
					Repo:  "x",
				},
				Path: "z",
				Ref:  "r",
			},
		},
		{
			input: map[string]any{"repo": "x/y", "path": "z"},
			want: github.File{
				Repo: github.Repo{
					Owner: github.User{Login: "x"},
					Repo:  "y",
				},
				Path: "z",
			},
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("%v", test.input), func(t *testing.T) {
			t.Parallel()

			got, err := parseFileMap(test.input)
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

	tests := []struct {
		input string
		want  github.File
	}{
		{
			input: "https://github.com/owner/repo/blob/ref/a",
			want: github.File{
				Repo: github.Repo{
					Owner: github.User{Login: "owner"},
					Repo:  "repo",
				},
				Path: "a",
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
			input: "owner/repo/blob/ref/a",
			want: github.File{
				Repo: github.Repo{
					Owner: github.User{Login: "owner"},
					Repo:  "repo",
				},
				Path: "a",
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
			input: "a@ref",
			want:  github.File{Path: "a", Ref: "ref"},
		},
		{
			input: complexPath + "@ref",
			want:  github.File{Path: complexPath, Ref: "ref"},
		},

		{
			input: "a",
			want:  github.File{Path: "a"},
		},
		{
			input: complexPath,
			want:  github.File{Path: complexPath},
		},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			t.Parallel()

			got, err := parseFileString(test.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if !test.want.Equal(got) {
				t.Errorf("want %+v, but got %+v", test.want, got)
			}
		})
	}
}
