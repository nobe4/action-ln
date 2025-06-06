package config

import (
	"testing"

	"github.com/nobe4/action-ln/internal/github"
)

func TestParseDefault(t *testing.T) {
	t.Parallel()

	t.Run("parses no link", func(t *testing.T) {
		t.Parallel()

		repo := github.Repo{}

		c := New(github.File{}, repo)
		raw := RawDefaults{}

		if err := c.parseDefaults(raw); err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if !c.Defaults.Link.From.Repo.Equal(repo) && c.Defaults.Link.To.Repo.Equal(repo) {
			t.Fatalf("expected default repo %v, got %v", repo, c.Defaults.Link)
		}
	})

	t.Run("parses one link", func(t *testing.T) {
		t.Parallel()

		c := New(github.File{}, github.Repo{})
		raw := RawDefaults{
			Link: RawLink{
				From: "o1/r1:p1",
				To:   "o2/r2:p2",
			},
		}

		want := Link{
			From: github.File{
				Repo: github.Repo{Owner: github.User{Login: "o1"}, Repo: "r1"},
				Path: "p1",
			},
			To: github.File{
				Repo: github.Repo{Owner: github.User{Login: "o2"}, Repo: "r2"},
				Path: "p2",
			},
		}

		if err := c.parseDefaults(raw); err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if !want.Equal(c.Defaults.Link) {
			t.Fatalf("expected %v, got %v", want, c.Defaults.Link)
		}
	})

	t.Run("parses more than one link", func(t *testing.T) {
		t.Parallel()

		c := New(github.File{}, github.Repo{})
		raw := RawDefaults{
			Link: RawLink{
				From: []any{
					"o1/r1:p1",
					"o2/r2:p2",
				},
				To: "o3/r3:p3",
			},
		}

		want := Link{
			From: github.File{
				Repo: github.Repo{Owner: github.User{Login: "o1"}, Repo: "r1"},
				Path: "p1",
			},
			To: github.File{
				Repo: github.Repo{Owner: github.User{Login: "o3"}, Repo: "r3"},
				Path: "p3",
			},
		}

		if err := c.parseDefaults(raw); err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if !want.Equal(c.Defaults.Link) {
			t.Fatalf("expected %v, got %v", want, c.Defaults.Link)
		}
	})
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

			c := New(github.File{}, github.Repo{})

			link, err := c.ParseLinkString(test.link)
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}

			want, err := c.ParseLinkString(test.want)
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}

			// force to only consider the From to To filling
			c.Defaults.Link = nil

			c.fillMissing(&link)

			if !link.Equal(&want) {
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

			c := New(github.File{}, github.Repo{})

			link, err := c.ParseLinkString(test.link)
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}

			want, err := c.ParseLinkString(test.want)
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}

			defaults, err := c.ParseLinkString(test.def)
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}

			c.Defaults.Link = &defaults

			c.fillDefaults(&link)

			if !link.Equal(&want) {
				t.Fatalf("expected\n%v => %v\ngot\n%v => %v", test.want, want, test.link, link)
			}
		})
	}
}
