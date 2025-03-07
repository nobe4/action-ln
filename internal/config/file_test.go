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

			// TODO: test err
			f, _ := parseFileMap(test.input)
			if !test.want.Equal(f) {
				t.Errorf("want %+v, but got %+v", test.want, f)
			}
		})
	}
}

func TestGetMapKey(t *testing.T) {
	t.Parallel()

	m := map[string]any{
		"a": "b",
		"c": 2,
	}

	if got := getMapKey(m, "a"); got != "b" {
		t.Errorf("want b, but got %v", got)
	}

	if got := getMapKey(m, "c"); got != "" {
		t.Errorf("want \"\", but got %v", got)
	}
}
