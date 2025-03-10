package config

import (
	"github.com/nobe4/action-ln/internal/github"
)

type Defaults struct {
	Repo github.Repo `json:"repo" yaml:"repo"`
}

func (d *Defaults) Equal(o *Defaults) bool {
	return d.Repo.Equal(o.Repo)
}

func (d *Defaults) parse(raw map[string]any) {
	if r := parseRepoString(
		getMapKey(raw, "owner"),
		getMapKey(raw, "repo"),
	); !r.Empty() {
		d.Repo = r
	}
}
