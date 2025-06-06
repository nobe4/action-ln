package config

import (
	"errors"
	"testing"

	fmock "github.com/nobe4/action-ln/internal/format/mock"
	"github.com/nobe4/action-ln/internal/github"
	gmock "github.com/nobe4/action-ln/internal/github/mock"
)

const (
	gotTo   = "got to"
	content = "content"
)

var errTest = errors.New("test")

func TestLinkNeedUpdate(t *testing.T) {
	t.Parallel()

	t.Run("content is the same on base branch", func(t *testing.T) {
		t.Parallel()

		g := gmock.Getter{
			FileHandler: func(_ *github.File) error { return errTest },
		}
		head := github.Branch{New: false}
		l := &Link{
			From: github.File{Content: content, Ref: "main"},
			To:   github.File{Content: content},
		}

		needUpdate, err := l.NeedUpdate(t.Context(), g, head)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if needUpdate {
			t.Fatalf("expected false, got %v", needUpdate)
		}
	})

	t.Run("to is missing", func(t *testing.T) {
		t.Parallel()

		g := gmock.Getter{
			FileHandler: func(_ *github.File) error { return github.ErrMissingFile },
		}
		head := github.Branch{New: false}
		l := &Link{
			From: github.File{Content: content, Ref: "main"},
			To:   github.File{Content: "content2"},
		}

		needUpdate, err := l.NeedUpdate(t.Context(), g, head)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if !needUpdate {
			t.Fatalf("expected true, got %v", needUpdate)
		}
	})

	t.Run("can't get the to file", func(t *testing.T) {
		t.Parallel()

		//nolint:err113 // This is just for this test.
		errWant := errors.New("test")

		g := gmock.Getter{
			FileHandler: func(_ *github.File) error { return errWant },
		}
		head := github.Branch{New: false}
		l := &Link{
			From: github.File{Content: content, Ref: "main"},
			To:   github.File{Content: "content2"},
		}

		_, err := l.NeedUpdate(t.Context(), g, head)
		if !errors.Is(err, errWant) {
			t.Fatalf("want error %v, got %v", errWant, err)
		}
	})

	t.Run("content is the same on head branch", func(t *testing.T) {
		t.Parallel()

		g := gmock.Getter{
			FileHandler: func(f *github.File) error {
				f.Content = content

				return nil
			},
		}
		head := github.Branch{New: false}
		l := &Link{
			From: github.File{Content: content, Ref: "main"},
			To:   github.File{Content: "content2"},
		}

		needUpdate, err := l.NeedUpdate(t.Context(), g, head)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if needUpdate {
			t.Fatalf("expected false, got %v", needUpdate)
		}
	})

	t.Run("content is different on head branch", func(t *testing.T) {
		t.Parallel()

		g := gmock.Getter{
			FileHandler: func(f *github.File) error {
				f.Content = "content2"

				return nil
			},
		}
		head := github.Branch{New: false}
		l := &Link{
			From: github.File{Content: content, Ref: "main"},
			To:   github.File{Content: "content2"},
		}

		needUpdate, err := l.NeedUpdate(t.Context(), g, head)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if !needUpdate {
			t.Fatalf("expected true, got %v", needUpdate)
		}
	})
}

func TestParseLinkString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		s    string
		want Link
		err  error
	}{
		{},

		{
			s: "a -> b",
			want: Link{
				From: github.File{Path: "a"},
				To:   github.File{Path: "b"},
			},
		},

		{
			s: "o/r:p@r -> b",
			want: Link{
				From: github.File{
					Repo: github.Repo{
						Owner: github.User{Login: "o"},
						Repo:  "r",
					},
					Path: "p",
					Ref:  "r",
				},
				To: github.File{Path: "b"},
			},
		},

		{
			s:   "a -> b -> c",
			err: errInvalidLinkFormat,
		},
	}

	for _, test := range tests {
		t.Run(test.s, func(t *testing.T) {
			t.Parallel()

			c := Config{}

			got, err := c.ParseLinkString(test.s)

			if !errors.Is(err, test.err) {
				t.Fatalf("expected error %v, got %v", test.err, err)
			}

			if !test.want.Equal(&got) {
				t.Fatalf("expected\n%+v\ngot\n%+v", test.want, got)
			}
		})
	}
}

func TestPopulate(t *testing.T) {
	t.Parallel()

	t.Run("fails to get from", func(t *testing.T) {
		t.Parallel()

		f := gmock.Getter{
			RepoHandler: func(_ *github.Repo) error { return errTest },
		}

		l := &Link{}

		if err := l.populate(t.Context(), f); !errors.Is(err, errGettingRepo) {
			t.Fatalf("expected error %v, got %v", errGettingRepo, err)
		}
	})

	t.Run("fails to get to", func(t *testing.T) {
		t.Parallel()

		f := gmock.Getter{
			FileHandler: func(f *github.File) error {
				if f.Path == "from" {
					return nil
				}

				return errTest
			},
		}

		l := &Link{
			From: github.File{Path: "from", Ref: "main"},
		}

		if err := l.populate(t.Context(), f); !errors.Is(err, errMissingTo) {
			t.Fatalf("expected error %v, got %v", errMissingTo, err)
		}
	})

	t.Run("succeeds with a missing to", func(t *testing.T) {
		t.Parallel()

		f := gmock.Getter{
			FileHandler: func(f *github.File) error {
				if f.Path == "from" {
					f.Content = "got"

					return nil
				}

				return github.ErrMissingFile
			},
		}

		l := &Link{
			From: github.File{Path: "from", Ref: "main"},
			To:   github.File{Path: "to"},
		}

		if err := l.populate(t.Context(), f); err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if l.From.Content != "got" {
			t.Fatalf("expected from to be populated, got %#v", l.From)
		}
	})

	t.Run("succeeds", func(t *testing.T) {
		t.Parallel()

		f := gmock.Getter{
			FileHandler: func(f *github.File) error {
				f.Content = "got " + f.Path

				return nil
			},
		}

		l := &Link{
			From: github.File{Path: "from", Ref: "main"},
			To:   github.File{Path: "to"},
		}

		if err := l.populate(t.Context(), f); err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if l.From.Content != "got from" {
			t.Fatalf("expected from to be populated, got %#v", l.From)
		}

		if l.To.Content != gotTo {
			t.Fatalf("expected from to be populated, got %#v", l.To)
		}
	})
}

func TestPopulateFrom(t *testing.T) {
	t.Parallel()

	const branch = "branch"

	t.Run("fails to get the repo", func(t *testing.T) {
		t.Parallel()

		f := gmock.Getter{
			RepoHandler: func(_ *github.Repo) error { return errTest },
		}

		l := &Link{}

		if err := l.populateFrom(t.Context(), f); !errors.Is(err, errGettingRepo) {
			t.Fatalf("expected error %v, got %v", errGettingRepo, err)
		}
	})

	t.Run("fails to get the file", func(t *testing.T) {
		t.Parallel()

		f := gmock.Getter{
			FileHandler: func(_ *github.File) error { return errTest },
			RepoHandler: func(r *github.Repo) error {
				r.DefaultBranch = branch

				return nil
			},
		}

		l := &Link{}

		if err := l.populateFrom(t.Context(), f); !errors.Is(err, errMissingFrom) {
			t.Fatalf("expected error %v, got %v", errMissingFrom, err)
		}
	})

	t.Run("gets the file with default ref", func(t *testing.T) {
		t.Parallel()

		f := gmock.Getter{
			FileHandler: func(f *github.File) error {
				f.Content = content

				return nil
			},
			RepoHandler: func(r *github.Repo) error {
				r.DefaultBranch = branch

				return nil
			},
		}

		l := &Link{}

		if err := l.populateFrom(t.Context(), f); err != nil {
			t.Fatalf("expected no error got %v", err)
		}

		if l.From.Content != content {
			t.Fatalf("expected content to be 'content', got %#v", l.From.Content)
		}

		if l.From.Ref != branch {
			t.Fatalf("expected ref to '%s', got %#v", branch, l.From.Ref)
		}
	})

	t.Run("gets the file with set ref", func(t *testing.T) {
		t.Parallel()

		f := gmock.Getter{
			FileHandler: func(f *github.File) error {
				f.Content = content

				return nil
			},
			RepoHandler: func(_ *github.Repo) error {
				t.Fatalf("RepoHandler should not be called in this test")

				return nil
			},
		}

		l := &Link{
			From: github.File{Ref: "main"},
		}

		if err := l.populateFrom(t.Context(), f); err != nil {
			t.Fatalf("expected no error got %v", err)
		}

		if l.From.Content != content {
			t.Fatalf("expected content to be 'content', got %#v", l.From.Content)
		}

		if l.From.Ref != "main" {
			t.Fatalf("expected ref to stay to 'main', got %#v", l.From.Ref)
		}
	})
}

func TestPopulateTo(t *testing.T) {
	t.Parallel()

	t.Run("fails to get the file on head branch", func(t *testing.T) {
		t.Parallel()

		f := gmock.Getter{
			FileHandler: func(_ *github.File) error { return errTest },
		}

		l := &Link{}

		if err := l.populateTo(t.Context(), f); !errors.Is(err, errMissingTo) {
			t.Fatalf("expected error %v, got %v", errMissingTo, err)
		}
	})

	t.Run("gets the file on head branch", func(t *testing.T) {
		t.Parallel()

		f := gmock.Getter{
			FileHandler: func(f *github.File) error {
				f.Content = gotTo

				return nil
			},
		}

		l := &Link{}

		if err := l.populateTo(t.Context(), f); err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if l.To.Content != gotTo {
			t.Fatalf("expected to be populated, got %#v", l.To)
		}
	})

	t.Run("fails to get the file on ref", func(t *testing.T) {
		t.Parallel()

		i := 0
		f := gmock.Getter{
			FileHandler: func(_ *github.File) error {
				if i == 0 {
					i++

					return github.ErrMissingFile
				}

				return errTest
			},
		}

		l := &Link{}

		if err := l.populateTo(t.Context(), f); !errors.Is(err, errMissingTo) {
			t.Fatalf("expected error %v, got %v", errMissingTo, err)
		}
	})

	t.Run("gets the file on ref", func(t *testing.T) {
		t.Parallel()

		i := 0
		f := gmock.Getter{
			FileHandler: func(f *github.File) error {
				if i == 0 {
					i++

					return github.ErrMissingFile
				}

				f.Content = gotTo

				return nil
			},
		}

		l := &Link{}

		if err := l.populateTo(t.Context(), f); err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if l.To.Content != gotTo {
			t.Fatalf("expected to be populated, got %#v", l.To)
		}
	})

	t.Run("both files are missing", func(t *testing.T) {
		t.Parallel()

		f := gmock.Getter{
			FileHandler: func(_ *github.File) error {
				return github.ErrMissingFile
			},
		}

		l := &Link{}

		if err := l.populateTo(t.Context(), f); err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if l.To.Content != "" {
			t.Fatalf("expected to not be populated, got %#v", l.To)
		}
	})
}

func TestLinkUpdate(t *testing.T) {
	t.Parallel()

	head := github.Branch{Name: "head"}

	t.Run("fail to update", func(t *testing.T) {
		t.Parallel()

		g := gmock.Updater{
			Handler: func(_ github.File, _ string, _ string) (github.File, error) {
				return github.File{}, errTest
			},
		}

		l := &Link{
			To:   github.File{Content: "to"},
			From: github.File{Content: "from"},
		}

		err := l.Update(t.Context(), g, fmock.New(), head)

		if !errors.Is(err, errTest) {
			t.Fatalf("want error %v, got %v", errTest, err)
		}

		if l.To.Content != l.From.Content {
			t.Fatal("want link content to be updated but isn't")
		}
	})

	t.Run("update", func(t *testing.T) {
		t.Parallel()

		g := gmock.Updater{
			Handler: func(_ github.File, _ string, _ string) (github.File, error) {
				// NOTE: the returned file is currently not used.
				return github.File{}, nil
			},
		}

		l := &Link{
			To:   github.File{Content: "to"},
			From: github.File{Content: "from"},
		}

		if err := l.Update(t.Context(), g, fmock.New(), head); err != nil {
			t.Fatalf("want no error, got %v", err)
		}

		if l.To.Content != l.From.Content {
			t.Fatal("want link content to be updated but isn't")
		}
	})
}

func TestParseLink(t *testing.T) {
	t.Parallel()

	repo := github.Repo{Owner: github.User{Login: "owner"}, Repo: "repo"}

	tests := []struct {
		rl   RawLink
		want Links
	}{
		{
			rl: RawLink{
				From: "from",
				To:   "to",
			},
			want: Links{
				{
					From: github.File{Path: "from"},
					To:   github.File{Path: "to"},
				},
			},
		},

		{
			rl: RawLink{From: "from", To: "to"},
			want: Links{
				{
					From: github.File{Path: "from"},
					To:   github.File{Path: "to"},
				},
			},
		},

		{
			rl: RawLink{
				From: map[string]any{"path": "from", "repo": "repo"},
				To:   "to",
			},
			want: Links{
				{
					From: github.File{Path: "from", Repo: github.Repo{Repo: "repo"}},
					To:   github.File{Path: "to", Repo: github.Repo{Repo: "repo"}},
				},
			},
		},

		{
			rl: RawLink{
				From: map[string]any{"path": "from", "repo": "repo", "owner": "owner"},
				To:   "to",
			},
			want: Links{
				{
					From: github.File{Path: "from", Repo: repo},
					To:   github.File{Path: "to", Repo: repo},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run("", func(t *testing.T) {
			t.Parallel()

			c := New(github.File{}, github.Repo{})

			got, err := c.parseLink(test.rl)
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}

			if !test.want.Equal(got) {
				t.Fatalf("expected\n%+v\ngot\n%+v", test.want, got)
			}
		})
	}
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

			link.fillMissing()

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

			link.fillDefaults(c.Defaults)

			if !link.Equal(&want) {
				t.Fatalf("expected\n%v => %v\ngot\n%v => %v", test.want, want, test.link, link)
			}
		})
	}
}
