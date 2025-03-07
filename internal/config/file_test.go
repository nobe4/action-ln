package config

import (
	"fmt"
	"testing"
)

func TestParseFileMap(t *testing.T) {
	t.Parallel()

	tests := []struct {
		input map[string]any
		want  File
	}{
		{},
		{
			input: map[string]any{"path": "z"},
			want:  File{Path: "z"},
		},
		{
			input: map[string]any{"repo": "x", "path": "z"},
			want:  File{Repo: "x", Path: "z"},
		},
		{
			input: map[string]any{"repo": "x", "owner": "y", "path": "z"},
			want:  File{Owner: "y", Repo: "x", Path: "z"},
		},
		{
			input: map[string]any{"repo": "x/y", "path": "z"},
			want:  File{Repo: "y", Owner: "x", Path: "z"},
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
		want  File
	}{
		{
			input: "https://github.com/owner/repo/blob/ref/a",
			want:  File{Owner: "owner", Repo: "repo", Path: "a"},
		},
		{
			input: "https://github.com/owner/repo/blob/ref/" + complexPath,
			want:  File{Owner: "owner", Repo: "repo", Path: complexPath},
		},

		{
			input: "owner/repo/blob/ref/a",
			want:  File{Owner: "owner", Repo: "repo", Path: "a"},
		},
		{
			input: "owner/repo/blob/ref/" + complexPath,
			want:  File{Owner: "owner", Repo: "repo", Path: complexPath},
		},

		{
			input: "a@ref",
			want:  File{Path: "a"},
		},
		{
			input: complexPath + "@ref",
			want:  File{Path: complexPath},
		},

		{
			input: "a",
			want:  File{Path: "a"},
		},
		{
			input: complexPath,
			want:  File{Path: complexPath},
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
