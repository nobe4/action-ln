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
	ErrNoPull     = errors.New("pull not found")
	ErrGetPull    = errors.New("failed to get pull")
	ErrCreatePull = errors.New("failed to create pull")
	ErrPullExists = errors.New("pull already exist")
)

type Pull struct {
	Number int `json:"number"`

	Repo Repo
	New  bool
}

func (p Pull) String() string {
	return fmt.Sprintf("https://github.com/%s/pull/%d", p.Repo, p.Number)
}

// https://docs.github.com/en/rest/pulls/pulls?apiVersion=2022-11-28#list-pull-requests
func (g *GitHub) GetPull(ctx context.Context, repo Repo, base, head string) (Pull, error) {
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

	path := fmt.Sprintf("/repos/%s/pulls?%s", repo, q.Encode())

	pulls := []Pull{}
	if _, err := g.req(ctx, http.MethodGet, path, nil, &pulls); err != nil {
		return Pull{}, fmt.Errorf("%w: %w", ErrGetPull, err)
	}

	if len(pulls) == 0 {
		return Pull{}, ErrNoPull
	}

	pull := pulls[0]
	pull.Repo = repo

	return pull, nil
}

// https://docs.github.com/en/rest/pulls/pulls?apiVersion=2022-11-28#create-a-pull-request
func (g *GitHub) CreatePull(ctx context.Context, repo Repo, base, head, title, pullBody string) (Pull, error) {
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
		return Pull{}, fmt.Errorf("%w: %w", ErrMarshalRequest, err)
	}

	path := fmt.Sprintf("/repos/%s/pulls", repo)

	pull := Pull{Repo: repo, New: true}
	if status, err := g.req(ctx, http.MethodPost, path, bytes.NewReader(body), &pull); err != nil {
		if status == http.StatusUnprocessableEntity {
			return Pull{}, ErrPullExists
		}

		return Pull{}, fmt.Errorf("%w: %w", ErrCreatePull, err)
	}

	return pull, nil
}

func (g *GitHub) GetOrCreatePull(ctx context.Context, repo Repo, base, head, title, body string) (Pull, error) {
	p, err := g.GetPull(ctx, repo, base, head)
	if err == nil {
		return p, nil
	}

	if !errors.Is(err, ErrNoPull) {
		return Pull{}, err
	}

	return g.CreatePull(ctx, repo, base, head, title, body)
}
