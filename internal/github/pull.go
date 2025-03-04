package github

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

var errNoPull = errors.New("no pull requests")

type Pull struct {
	Number int `json:"number"`
}

// https://docs.github.com/en/rest/pulls/pulls?apiVersion=2022-11-28#list-pull-requests
func (g GitHub) GetPull(ctx context.Context, repo Repo, base, head string) (Pull, error) {
	q := url.Values{
		"base": []string{base},
		"head": []string{repo.Owner.Login + ":" + head},

		// NOTE: GitHub only ever allows 1 PR per HEAD/BASE branches.
		// If you try to create a PR with the same branches, it will fail with:
		// {
		//   "status": "422"
		//   "errors": [ { "message": "A pull request already exists for <OWNER>:<HEAD>." } ],
		// }
		"per_page": []string{"1"},

		"state": []string{"open"},
	}

	path := fmt.Sprintf("/repos/%s/%s/pulls?%s", repo.Owner.Login, repo.Repo, q.Encode())

	pulls := []Pull{}
	if _, err := g.req(ctx, http.MethodGet, path, nil, &pulls); err != nil {
		return Pull{}, fmt.Errorf("failed to get pulls: %w", err)
	}

	if len(pulls) == 0 {
		return Pull{}, errNoPull
	}

	return pulls[0], nil
}
