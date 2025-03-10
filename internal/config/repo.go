package config

import (
	"strings"

	"github.com/nobe4/action-ln/internal/github"
)

func parseRepoString(owner, repo string) github.Repo {
	r := github.Repo{
		Owner: github.User{Login: owner},
		Repo:  repo,
	}

	if strings.Contains(repo, "/") {
		parts := strings.Split(repo, "/")
		if len(parts) != 2 { //nolint:all // TODO: log that the repo is badly formatted
			// don't do anything
		}

		r.Owner.Login = parts[0]
		r.Repo = parts[1]
	}

	return r
}
