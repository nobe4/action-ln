package config

import (
	"testing"

	fmock "github.com/nobe4/action-ln/internal/format/mock"
	"github.com/nobe4/action-ln/internal/github"
	gmock "github.com/nobe4/action-ln/internal/github/mock"
)

func TestCombineLinks(t *testing.T) {
	t.Parallel()

	mkFile := func(n string) github.File {
		return github.File{
			Path: n,
		}
	}
	mkLink := func(from, to string) *Link {
		return &Link{
			From: mkFile(from),
			To:   mkFile(to),
		}
	}

	tests := []struct {
		froms []github.File
		tos   []github.File
		want  Links
	}{
		{},

		{
			froms: []github.File{},
			tos: []github.File{
				mkFile("0"),
				mkFile("1"),
			},
			want: Links{},
		},

		{
			froms: []github.File{
				mkFile("0"),
				mkFile("1"),
			},
			tos: []github.File{},
			want: Links{
				mkLink("0", "0"),
				mkLink("1", "1"),
			},
		},

		{
			froms: []github.File{
				mkFile("0"),
				mkFile("1"),
			},
			tos: []github.File{
				mkFile("2"),
			},
			want: Links{
				mkLink("0", "2"),
				mkLink("1", "2"),
			},
		},

		{
			froms: []github.File{
				mkFile("0"),
				mkFile("1"),
			},
			tos: []github.File{
				mkFile("2"),
				mkFile("3"),
			},
			want: Links{
				mkLink("0", "2"),
				mkLink("0", "3"),
				mkLink("1", "2"),
				mkLink("1", "3"),
			},
		},
	}

	for _, test := range tests {
		t.Run("", func(t *testing.T) {
			t.Parallel()

			got := combineLinks(test.froms, test.tos)
			if !got.Equal(test.want) {
				t.Fatalf("expected %+v, got %+v", test.want, got)
			}
		})
	}
}

func TestFilter(t *testing.T) {
	t.Parallel()

	mkLink := func(f, t string) *Link {
		return &Link{
			From: github.File{
				Path: f + "_path",
				Repo: github.Repo{
					Owner: github.User{Login: f + "_owner"},
					Repo:  f + "_repo",
				},
			},
			To: github.File{
				Path: t + "_path",
				Repo: github.Repo{
					Owner: github.User{Login: t + "_owner"},
					Repo:  t + "_repo",
				},
			},
		}
	}

	tests := []struct {
		links Links
		want  Links
	}{
		{},

		{
			links: Links{
				mkLink("a", "b"),
				mkLink("c", "d"),
				mkLink("e", "f"),
			},
			want: Links{
				mkLink("a", "b"),
				mkLink("c", "d"),
				mkLink("e", "f"),
			},
		},

		{
			links: Links{
				mkLink("a", "b"),
				mkLink("c", "c"),
				mkLink("d", "e"),
			},
			want: Links{
				mkLink("a", "b"),
				mkLink("d", "e"),
			},
		},

		{
			links: Links{
				mkLink("a", "a"),
				mkLink("b", "b"),
				mkLink("c", "c"),
			},
			want: Links{},
		},
	}

	for _, test := range tests {
		t.Run("", func(t *testing.T) {
			t.Parallel()

			test.links.Filter()

			if !test.want.Equal(test.links) {
				t.Fatalf("expected %+v, got %+v", test.want, test.links)
			}
		})
	}
}

func TestLinksUpdate(t *testing.T) {
	t.Parallel()

	head := github.Branch{New: false}

	const got = "got"

	t.Run("fail to check if the link needs an update", func(t *testing.T) {
		t.Parallel()

		g := gmock.FileGetterUpdater{
			GetHandler: func(*github.File) error { return errTest },
		}

		l := &Links{
			{
				From: github.File{Content: "from"},
				To:   github.File{Content: "to"},
			},
		}

		updated := l.Update(t.Context(), g, fmock.New(), head)

		if s := (*l)[0].Status; s != "failed to check for update" {
			t.Fatalf("want status 'failed to check for update', got '%s'", s)
		}

		if updated {
			t.Fatal("want to not be updated")
		}
	})

	t.Run("do not update the link", func(t *testing.T) {
		t.Parallel()

		g := gmock.FileGetterUpdater{}

		l := &Links{
			{
				From: github.File{Content: "from"},
				To:   github.File{Content: "from"},
			},
		}

		updated := l.Update(t.Context(), g, fmock.New(), head)
		if updated {
			t.Fatal("want to not be updated")
		}
	})

	t.Run("fail to update the link", func(t *testing.T) {
		t.Parallel()

		g := gmock.FileGetterUpdater{
			GetHandler: func(f *github.File) error {
				f.Content = got

				return nil
			},
			UpdateHandler: func(github.File, string, string) (github.File, error) {
				return github.File{}, errTest
			},
		}

		l := &Links{
			{
				From: github.File{Content: "from"},
				To:   github.File{Content: "to"},
			},
		}

		updated := l.Update(t.Context(), g, fmock.New(), head)

		if s := (*l)[0].Status; s != "failed to update" {
			t.Fatalf("want status 'failed to update', got '%s'", s)
		}

		if updated {
			t.Fatal("want to not be updated")
		}
	})

	t.Run("update the link", func(t *testing.T) {
		t.Parallel()

		g := gmock.FileGetterUpdater{
			GetHandler: func(f *github.File) error {
				f.Content = got

				return nil
			},
			UpdateHandler: func(github.File, string, string) (github.File, error) {
				return github.File{}, nil
			},
		}

		l := &Links{
			{
				From: github.File{Content: "from"},
				To:   github.File{Content: "to"},
			},
		}

		updated := l.Update(t.Context(), g, fmock.New(), head)

		if !updated {
			t.Fatal("want to be updated")
		}
	})

	t.Run("multiple links", func(t *testing.T) {
		t.Parallel()

		g := gmock.FileGetterUpdater{
			GetHandler: func(f *github.File) error {
				f.Content = got

				return nil
			},
			UpdateHandler: func(f github.File, _ string, _ string) (github.File, error) {
				if f.Content == "error" {
					return github.File{}, errTest
				}

				return github.File{}, nil
			},
		}

		l := &Links{
			// Needs no update
			{
				From: github.File{Content: "from"},
				To:   github.File{Content: "from"},
			},

			// Updates correctly
			{
				From: github.File{Content: "from"},
				To:   github.File{Content: "to"},
			},

			// Fails to update
			{
				From: github.File{Content: "error"},
				To:   github.File{Content: "to"},
			},
		}

		updated := l.Update(t.Context(), g, fmock.New(), head)

		if s := (*l)[0].Status; s != "update not needed" {
			t.Fatalf("want status 'update not needed', got '%s'", s)
		}

		if s := (*l)[1].Status; s != "updated" {
			t.Fatalf("want status 'updated', got '%s'", s)
		}

		if s := (*l)[2].Status; s != "failed to update" {
			t.Fatalf("want status 'failed to update', got '%s'", s)
		}

		// The function failed but we had a valid update.
		if !updated {
			t.Fatal("want to be updated")
		}
	})
}

func TestGroups(t *testing.T) {
	t.Parallel()

	links := Links{
		&Link{
			To: github.File{
				Repo: github.Repo{Owner: github.User{Login: "a"}, Repo: "b"},
			},
		},

		&Link{
			To: github.File{
				Repo: github.Repo{Owner: github.User{Login: "a"}, Repo: "b"},
			},
		},

		&Link{
			To: github.File{
				Repo: github.Repo{Owner: github.User{Login: "a"}, Repo: "c"},
			},
		},

		&Link{
			To: github.File{
				Repo: github.Repo{Owner: github.User{Login: "d"}, Repo: "e"},
			},
		},
	}

	got := links.Groups()

	if got["a/b"][0] != links[0] {
		t.Fatalf("expected %v, got %v", links[0], got["a/b"][0])
	}

	if got["a/b"][1] != links[1] {
		t.Fatalf("expected %v, got %v", links[1], got["a/b"][1])
	}

	if got["a/c"][0] != links[2] {
		t.Fatalf("expected %v, got %v", links[2], got["a/c"][0])
	}

	if got["d/e"][0] != links[3] {
		t.Fatalf("expected %v, got %v", links[3], got["d/e"][0])
	}
}
