package config

import (
	"testing"

	"github.com/nobe4/action-ln/internal/github"
)

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
