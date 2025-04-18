package config

import (
	"strings"

	"github.com/nobe4/action-ln/internal/github"
	"github.com/nobe4/action-ln/internal/log"
)

const repoPartsCount = 2 // owner/repo

func parseRepoString(owner, repo string) github.Repo {
	r := github.Repo{
		Owner: github.User{Login: owner},
		Repo:  repo,
	}

	if strings.Contains(repo, "/") {
		parts := strings.Split(repo, "/")
		if len(parts) != repoPartsCount {
			log.Warn("Invalid repo string", "repo", repo)
		}

		r.Owner.Login = parts[0]
		r.Repo = parts[1]
	}

	log.Debug("Parse repo", "owner", owner, "repo", repo, "parsed", r)

	return r
}
