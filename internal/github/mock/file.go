package mock

import (
	"context"

	"github.com/nobe4/action-ln/internal/github"
)

type FileGetter struct {
	Handler func(*github.File) error
}

func (f FileGetter) GetFile(_ context.Context, file *github.File) error {
	return f.Handler(file)
}
