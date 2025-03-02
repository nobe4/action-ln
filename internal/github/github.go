/*
Package github implements common interactions with GitHub's API.

Refs:
- https://docs.github.com/en/rest/authentication/authenticating-to-the-rest-api?apiVersion=2022-11-28
*/
package github

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

var errRequest = errors.New("request failed")

const (
	PathUser = "/user"
)

type GitHub struct {
	client   http.Client
	token    string
	endpoint string
}

func New(token, endpoint string) GitHub {
	return GitHub{
		client:   http.Client{},
		token:    token,
		endpoint: endpoint,
	}
}

type User struct {
	Login string `json:"login"`
}

func (g GitHub) GetUser(ctx context.Context) (User, error) {
	u := User{}

	if err := g.req(ctx, "GET", PathUser, nil, &u); err != nil {
		return u, fmt.Errorf("failed to get user: %w", err)
	}

	return u, nil
}

//nolint:unparam // Will add more requests later
func (g GitHub) req(ctx context.Context, method, path string, body io.Reader, out any) error {
	path = g.endpoint + path

	req, err := http.NewRequestWithContext(ctx, method, path, body)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("Authorization", "Bearer "+g.token)

	res, err := g.client.Do(req)
	if err != nil {
		return fmt.Errorf("%w: %w", errRequest, err)
	}
	defer res.Body.Close()

	// All the 2XX codes
	success := res.StatusCode >= http.StatusOK && res.StatusCode < http.StatusMultipleChoices
	if !success {
		return fmt.Errorf("%w: %s", errRequest, res.Status)
	}

	if out != nil {
		if err := json.NewDecoder(res.Body).Decode(out); err != nil {
			return fmt.Errorf("failed to decode response: %w", err)
		}
	}

	return nil
}
