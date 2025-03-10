package config

import "github.com/nobe4/action-ln/internal/github"

type Defaults struct {
	Repo github.Repo `json:"repo" yaml:"repo"`
}

func (d Defaults) Equal(o Defaults) bool {
	return d.Repo.Equal(o.Repo)
}

func parseDefaults(raw map[string]any) Defaults {
	d := Defaults{}

	d.Repo = parseRepoString(
		getMapKey(raw, "owner"),
		getMapKey(raw, "repo"),
	)

	return d
}
