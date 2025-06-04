package mock

import (
	"context"

	"github.com/nobe4/action-ln/internal/github"
)

type Getter struct {
	FileHandler func(*github.File) error
	RepoHandler func(*github.Repo) error
}

func (g Getter) GetFile(_ context.Context, f *github.File) error {
	return g.FileHandler(f)
}

func (g Getter) GetRepo(_ context.Context, r *github.Repo) error {
	return g.RepoHandler(r)
}

type Updater struct {
	Handler func(github.File, string, string) (github.File, error)
}

func (g Updater) UpdateFile(_ context.Context, f github.File, head, msg string) (github.File, error) {
	return g.Handler(f, head, msg)
}

type GetterUpdater struct {
	GetFileHandler func(*github.File) error
	GetRepoHandler func(*github.Repo) error
	UpdateHandler  func(github.File, string, string) (github.File, error)
}

func (g GetterUpdater) GetFile(_ context.Context, f *github.File) error {
	return g.GetFileHandler(f)
}

func (g GetterUpdater) GetRepo(_ context.Context, r *github.Repo) error {
	return g.GetRepoHandler(r)
}

func (g GetterUpdater) UpdateFile(_ context.Context, f github.File, head, msg string) (github.File, error) {
	return g.UpdateHandler(f, head, msg)
}
