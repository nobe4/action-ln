package config

import (
	"fmt"
	"testing"

	"github.com/nobe4/action-ln/internal/github"
)

func TestParseDefaults(t *testing.T) {
	t.Parallel()

	tests := []struct {
		input map[string]any
		want  Defaults
	}{
		{},
		{
			input: map[string]any{"repo": "x"},
			want:  Defaults{Repo: github.Repo{Repo: "x"}},
		},
		{
			input: map[string]any{"repo": "x", "owner": "y"},
			want: Defaults{
				Repo: github.Repo{
					Owner: github.User{Login: "y"},
					Repo:  "x",
				},
			},
		},
		{
			input: map[string]any{"owner": "y"},
			want: Defaults{
				Repo: github.Repo{
					Owner: github.User{Login: "y"},
				},
			},
		},
		{
			input: map[string]any{"repo": "x/y"},
			want: Defaults{
				Repo: github.Repo{
					Owner: github.User{Login: "x"},
					Repo:  "y",
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("%v", test.input), func(t *testing.T) {
			t.Parallel()

			got := parseDefaults(test.input)

			if !test.want.Equal(got) {
				t.Errorf("want %+v, but got %+v", test.want, got)
			}
		})
	}
}
