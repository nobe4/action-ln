package mock

import (
	"context"

	"github.com/nobe4/action-ln/internal/github"
)

type FileGetter struct {
	Handler func(*github.File) error
}

func (g FileGetter) GetFile(_ context.Context, f *github.File) error {
	return g.Handler(f)
}

type FileRepoGetter struct {
	FileHandler func(*github.File) error
	RepoHandler func(*github.Repo) error
}

func (g FileRepoGetter) GetFile(_ context.Context, f *github.File) error {
	return g.FileHandler(f)
}

func (g FileRepoGetter) GetRepo(_ context.Context, r *github.Repo) error {
	return g.RepoHandler(r)
}

type FileUpdater struct {
	Handler func(github.File, string, string) (github.File, error)
}

func (g FileUpdater) UpdateFile(_ context.Context, f github.File, head, msg string) (github.File, error) {
	return g.Handler(f, head, msg)
}

type FileGetterUpdater struct {
	GetHandler    func(*github.File) error
	UpdateHandler func(github.File, string, string) (github.File, error)
}

func (g FileGetterUpdater) GetFile(_ context.Context, f *github.File) error {
	return g.GetHandler(f)
}

func (g FileGetterUpdater) UpdateFile(_ context.Context, f github.File, head, msg string) (github.File, error) {
	return g.UpdateHandler(f, head, msg)
}
