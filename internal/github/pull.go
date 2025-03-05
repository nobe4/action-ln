package github

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

var (
	errNoPull     = errors.New("no pull requests")
	errPullExists = errors.New("pull request already exists")
)

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

// https://docs.github.com/en/rest/pulls/pulls?apiVersion=2022-11-28#create-a-pull-request
func (g GitHub) CreatePull(ctx context.Context, repo Repo, base, head, title, pullBody string) (Pull, error) {
	body, err := json.Marshal(struct {
		Title string `json:"title"`
		Head  string `json:"head"`
		Base  string `json:"base"`
		Body  string `json:"body"`
	}{
		Title: title,
		Body:  pullBody,
		Head:  repo.Owner.Login + ":" + head,
		Base:  base,
	})
	if err != nil {
		return Pull{}, fmt.Errorf("failed to marshal request: %w", err)
	}

	path := fmt.Sprintf("/repos/%s/%s/pulls", repo.Owner.Login, repo.Repo)

	pull := Pull{}
	if status, err := g.req(ctx, http.MethodPost, path, bytes.NewReader(body), &pull); err != nil {
		if status == http.StatusUnprocessableEntity {
			return Pull{}, errPullExists
		}

		return Pull{}, fmt.Errorf("failed to create pull: %w", err)
	}

	return pull, nil
}
