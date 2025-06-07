package path

import (
	"slices"
	"testing"
)

func TestTrimN(t *testing.T) {
	t.Parallel()

	tests := []struct {
		p    string
		n    int
		want string
	}{
		{},

		{
			p:    "a",
			want: "a",
		},

		{
			p:    "a/b/c",
			want: "a/b/c",
		},

		{
			p:    "a/b/c",
			n:    0,
			want: "a/b/c",
		},

		{
			p:    "a/b/c",
			n:    1,
			want: "b/c",
		},

		{
			p:    "a/b/c",
			n:    -1,
			want: "a/b",
		},

		{
			p:    "a/b/c/d/e",
			n:    4,
			want: "e",
		},

		{
			p:    "a/b/c/d/e",
			n:    -3,
			want: "a/b",
		},
	}

	for _, test := range tests {
		t.Run("", func(t *testing.T) {
			t.Parallel()

			got := TrimN(test.p, test.n)

			if test.want != got {
				t.Errorf("want %q, but got %q", test.want, got)
			}
		})
	}
}

func TestSplitAll(t *testing.T) {
	t.Parallel()

	tests := []struct {
		p    string
		want []string
	}{
		{
			p:    "",
			want: []string{""},
		},

		{
			p:    "a",
			want: []string{"a"},
		},

		{
			p:    "a/b/c",
			want: []string{"a", "b", "c"},
		},

		{
			p:    "a/b/c/d.txt",
			want: []string{"a", "b", "c", "d.txt"},
		},
	}

	for _, test := range tests {
		t.Run("", func(t *testing.T) {
			t.Parallel()

			got := splitAll(test.p)

			if !slices.Equal(got, test.want) {
				t.Errorf("want %q, but got %q", test.want, got)
			}
		})
	}
}
