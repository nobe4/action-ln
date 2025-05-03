package config

import (
	"errors"
	"testing"

	fmock "github.com/nobe4/action-ln/internal/format/mock"
	"github.com/nobe4/action-ln/internal/github"
	gmock "github.com/nobe4/action-ln/internal/github/mock"
)

var errTest = errors.New("test")

func TestLinkNeedUpdate(t *testing.T) {
	t.Parallel()

	t.Run("content is the same on base branch", func(t *testing.T) {
		t.Parallel()

		g := gmock.FileGetter{Handler: func(_ *github.File) error { return errTest }}
		head := github.Branch{New: false}
		l := &Link{
			From: github.File{Content: "content"},
			To:   github.File{Content: "content"},
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

		g := gmock.FileGetter{
			Handler: func(_ *github.File) error { return github.ErrMissingFile },
		}
		head := github.Branch{New: false}
		l := &Link{
			From: github.File{Content: "content"},
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

		g := gmock.FileGetter{
			Handler: func(_ *github.File) error { return errWant },
		}
		head := github.Branch{New: false}
		l := &Link{
			From: github.File{Content: "content"},
			To:   github.File{Content: "content2"},
		}

		_, err := l.NeedUpdate(t.Context(), g, head)
		if !errors.Is(err, errWant) {
			t.Fatalf("want error %v, got %v", errWant, err)
		}
	})

	t.Run("content is the same on head branch", func(t *testing.T) {
		t.Parallel()

		g := gmock.FileGetter{
			Handler: func(f *github.File) error {
				f.Content = "content"

				return nil
			},
		}
		head := github.Branch{New: false}
		l := &Link{
			From: github.File{Content: "content"},
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

		g := gmock.FileGetter{
			Handler: func(f *github.File) error {
				f.Content = "content2"

				return nil
			},
		}
		head := github.Branch{New: false}
		l := &Link{
			From: github.File{Content: "content"},
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

func TestPopulate(t *testing.T) {
	t.Parallel()

	t.Run("fails to get from", func(t *testing.T) {
		t.Parallel()

		f := gmock.FileGetter{Handler: func(_ *github.File) error { return errTest }}

		l := &Link{}

		if err := l.populate(t.Context(), f); !errors.Is(err, errMissingFrom) {
			t.Fatalf("expected error %v, got %v", errMissingFrom, err)
		}
	})

	t.Run("fails to get to", func(t *testing.T) {
		t.Parallel()

		f := gmock.FileGetter{
			Handler: func(f *github.File) error {
				if f.Path == "from" {
					return nil
				}

				return errTest
			},
		}

		l := &Link{
			From: github.File{Path: "from"},
		}

		if err := l.populate(t.Context(), f); !errors.Is(err, errMissingTo) {
			t.Fatalf("expected error %v, got %v", errMissingTo, err)
		}
	})

	t.Run("succeeds with a missing to", func(t *testing.T) {
		t.Parallel()

		f := gmock.FileGetter{
			Handler: func(f *github.File) error {
				if f.Path == "from" {
					f.Content = "got"

					return nil
				}

				return github.ErrMissingFile
			},
		}

		l := &Link{
			From: github.File{Path: "from"},
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

		f := gmock.FileGetter{
			Handler: func(f *github.File) error {
				f.Content = "got " + f.Path

				return nil
			},
		}

		l := &Link{
			From: github.File{Path: "from"},
			To:   github.File{Path: "to"},
		}

		if err := l.populate(t.Context(), f); err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if l.From.Content != "got from" {
			t.Fatalf("expected from to be populated, got %#v", l.From)
		}

		if l.To.Content != "got to" {
			t.Fatalf("expected from to be populated, got %#v", l.To)
		}
	})
}

func TestLinkUpdate(t *testing.T) {
	t.Parallel()

	head := github.Branch{Name: "head"}

	t.Run("fail to update", func(t *testing.T) {
		t.Parallel()

		g := gmock.FileUpdater{
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

		g := gmock.FileUpdater{
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
