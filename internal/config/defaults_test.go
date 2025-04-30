package config

import (
	"strings"
	"testing"
)

// TODO: redo
// func TestDefaultsParse(t *testing.T) {
// 	t.Parallel()
//
// 	tests := []struct {
// 		input map[string]any
// 		want  Defaults
// 	}{
// 		{},
// 		{
// 			input: map[string]any{"repo": "x"},
// 			want:  Defaults{Repo: github.Repo{Repo: "x"}},
// 		},
// 		{
// 			input: map[string]any{"repo": "x", "owner": "y"},
// 			want: Defaults{
// 				Repo: github.Repo{
// 					Owner: github.User{Login: "y"},
// 					Repo:  "x",
// 				},
// 			},
// 		},
// 		{
// 			input: map[string]any{"owner": "y"},
// 			want: Defaults{
// 				Repo: github.Repo{
// 					Owner: github.User{Login: "y"},
// 				},
// 			},
// 		},
// 		{
// 			input: map[string]any{"repo": "x/y"},
// 			want: Defaults{
// 				Repo: github.Repo{
// 					Owner: github.User{Login: "x"},
// 					Repo:  "y",
// 				},
// 			},
// 		},
// 	}
//
// 	for _, test := range tests {
// 		t.Run(fmt.Sprintf("%v", test.input), func(t *testing.T) {
// 			t.Parallel()
//
// 			d := &Defaults{}
// 			d.parse(test.input)
//
// 			if !test.want.Equal(d) {
// 				t.Errorf("want %+v, but got %+v", test.want, d)
// 			}
// 		})
// 	}
//
// 	t.Run("overwite existing values", func(t *testing.T) {
// 		t.Parallel()
//
// 		input := map[string]any{
// 			"repo": "a/b",
// 		}
//
// 		want := &Defaults{
// 			Repo: github.Repo{
// 				Owner: github.User{Login: "a"},
// 				Repo:  "b",
// 			},
// 		}
//
// 		d := &Defaults{
// 			Repo: github.Repo{
// 				Owner: github.User{Login: "x"},
// 				Repo:  "y",
// 			},
// 		}
//
// 		d.parse(input)
//
// 		if !d.Equal(want) {
// 			t.Errorf("want %+v, but got %+v", want, d)
// 		}
// 	})
// }

// Parse link from `o/r:p -> o/r:p` for simpler test cases.
func parseLink(t *testing.T, c *Config, s string) *Link {
	t.Helper()

	if s == "" {
		return &Link{}
	}

	p := strings.Split(s, " -> ")

	if l := len(p); l != 2 {
		t.Fatalf("expected 2 parts, got %d for %q", l, s)
	}

	from, err := c.parseString(p[0])
	if err != nil {
		t.Fatalf("failed to parse file %q: %v", p[0], err)
	}

	to, err := c.parseString(p[1])
	if err != nil {
		t.Fatalf("failed to parse file %q: %v", p[0], err)
	}

	return &Link{From: from[0], To: to[0]}
}

func TestFillMissing(t *testing.T) {
	t.Parallel()

	tests := []struct {
		link string
		want string
	}{
		{},

		{
			link: "o1/r1:p1 -> o2/r2:p2",
			want: "o1/r1:p1 -> o2/r2:p2",
		},

		{
			link: "o1/r1:p1 -> p2",
			want: "o1/r1:p1 -> o1/r1:p2",
		},

		{
			link: "o1/:p1 -> p2",
			want: "o1/:p1 -> o1/:p2",
		},

		{
			link: "/r1:p1 -> p2",
			want: "/r1:p1 -> /r1:p2",
		},

		{
			link: "p1 -> o2/r2:p2",
			want: "p1 -> o2/r2:p2",
		},

		{
			link: "p1 -> o2/:p2",
			want: "p1 -> o2/:p2",
		},

		{
			link: "p1 -> /r2:p2",
			want: "p1 -> /r2:p2",
		},

		{
			link: "p1 -> p2",
			want: "p1 -> p2",
		},
	}

	for _, test := range tests {
		t.Run("", func(t *testing.T) {
			t.Parallel()

			c := New()

			link, want := parseLink(t, c, test.link), parseLink(t, c, test.want)

			// force to only consider the From to To filling
			c.Defaults.Link = nil

			c.fillMissing(link)

			if !link.Equal(want) {
				t.Fatalf("expected\n%v => %v\ngot\n%v => %v", test.want, want, test.link, link)
			}
		})
	}
}

func TestFillDefaults(t *testing.T) {
	t.Parallel()

	tests := []struct {
		def  string
		link string
		want string
	}{
		{},

		{
			link: "o1/r1:p1 -> o2/r2:p2",
			want: "o1/r1:p1 -> o2/r2:p2",
		},

		{
			link: "p1 -> o2/r2:p2",
			want: "p1 -> o2/r2:p2",
		},

		{
			link: "p1 -> o2/:p2",
			want: "p1 -> o2/:p2",
		},

		{
			link: "p1 -> /r2:p2",
			want: "p1 -> /r2:p2",
		},

		{
			link: "p1 -> p2",
			want: "p1 -> p2",
		},

		{
			def:  "o3/r3:p3 -> o4/r4:p4",
			link: "",
			want: "o3/r3:p3 -> o4/r4:p4",
		},

		{
			def:  "o3/r3:p3 -> o4/r4:p4",
			link: "o1/r1:p1 -> o2/r2:p2",
			want: "o1/r1:p1 -> o2/r2:p2",
		},

		{
			def:  "o3/r3:p3 -> o4/r4:p4",
			link: "o1/:p1 -> o2/:p2",
			want: "o1/:p1 -> o2/:p2",
		},

		{
			def:  "o3/r3:p3 -> o4/r4:p4",
			link: "/r1:p1 -> /r2:p2",
			want: "/r1:p1 -> /r2:p2",
		},

		{
			def:  "o3/r3:p3 -> o4/r4:p4",
			link: "o1/r1:p1 -> p2",
			want: "o1/r1:p1 -> o4/r4:p2",
		},

		{
			def:  "o3/r3:p3 -> o4/r4:p4",
			link: "p1 -> o2/r2:p2",
			want: "o3/r3:p1 -> o2/r2:p2",
		},

		{
			def:  "o3/r3:p3 -> o4/r4:p4",
			link: "p1 -> p2",
			want: "o3/r3:p1 -> o4/r4:p2",
		},
	}

	for _, test := range tests {
		t.Run("", func(t *testing.T) {
			t.Parallel()

			c := New()

			link, want, defaults := parseLink(t, c, test.link), parseLink(t, c, test.want), parseLink(t, c, test.def)

			c.Defaults.Link = defaults

			c.fillDefaults(link)

			if !link.Equal(want) {
				t.Fatalf("expected\n%v => %v\ngot\n%v => %v", test.want, want, test.link, link)
			}
		})
	}
}
