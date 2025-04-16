package config

import (
	"errors"
	"testing"

	"github.com/nobe4/action-ln/internal/github"
	"github.com/nobe4/action-ln/internal/github/mock"
)

func TestLinksUpdate(t *testing.T) {
	t.Parallel()

	head := github.Branch{New: false}

	const got = "got"

	t.Run("fail to check if the link needs an update", func(t *testing.T) {
		t.Parallel()

		g := mock.FileGetterUpdater{
			GetHandler: func(*github.File) error { return errTest },
		}

		l := &Links{
			{
				From: github.File{Content: "from"},
				To:   github.File{Content: "to"},
			},
		}

		updated, err := l.Update(t.Context(), g, head)
		if !errors.Is(err, errTest) {
			t.Fatalf("want error %v, got %v", errTest, err)
		}

		if updated {
			t.Fatal("want to not be updated")
		}
	})

	t.Run("do not update the link", func(t *testing.T) {
		t.Parallel()

		g := mock.FileGetterUpdater{}

		l := &Links{
			{
				From: github.File{Content: "from"},
				To:   github.File{Content: "from"},
			},
		}

		updated, err := l.Update(t.Context(), g, head)
		if err != nil {
			t.Fatalf("want no error, got %v", err)
		}

		if updated {
			t.Fatal("want to not be updated")
		}
	})

	t.Run("fail to update the link", func(t *testing.T) {
		t.Parallel()

		g := mock.FileGetterUpdater{
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

		updated, err := l.Update(t.Context(), g, head)
		if !errors.Is(err, errTest) {
			t.Fatalf("want error %v, got %v", errTest, err)
		}

		if updated {
			t.Fatal("want to not be updated")
		}
	})

	t.Run("update the link", func(t *testing.T) {
		t.Parallel()

		g := mock.FileGetterUpdater{
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

		updated, err := l.Update(t.Context(), g, head)
		if err != nil {
			t.Fatalf("want no error, got %v", err)
		}

		if !updated {
			t.Fatal("want to be updated")
		}
	})

	t.Run("multiple links", func(t *testing.T) {
		t.Parallel()

		g := mock.FileGetterUpdater{
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

		updated, err := l.Update(t.Context(), g, head)
		if !errors.Is(err, errTest) {
			t.Fatalf("want error %v, got %v", errTest, err)
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
