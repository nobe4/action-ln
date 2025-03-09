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

var (
	ErrRequestFailed  = errors.New("request failed")
	ErrMarshalRequest = errors.New("failed to marshal request")
)

const (
	PathUser = "/user"
)

type GitHub struct {
	client   http.Client
	Token    string
	endpoint string
}

func New(endpoint string) GitHub {
	return GitHub{
		client:   http.Client{},
		endpoint: endpoint,
	}
}

type User struct {
	Login string `json:"login"`
}

func (g *GitHub) req(ctx context.Context, method, path string, body io.Reader, out any) (int, error) {
	path = g.endpoint + path

	req, err := http.NewRequestWithContext(ctx, method, path, body)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("Authorization", "Bearer "+g.Token)

	res, err := g.client.Do(req)
	if err != nil {
		return res.StatusCode, fmt.Errorf("%w: %w", ErrRequestFailed, err)
	}
	defer res.Body.Close()

	// All the 2XX codes
	success := res.StatusCode >= http.StatusOK && res.StatusCode < http.StatusMultipleChoices
	if !success {
		return res.StatusCode, fmt.Errorf("%w: %s", ErrRequestFailed, res.Status)
	}

	if out != nil {
		if err := json.NewDecoder(res.Body).Decode(out); err != nil {
			return http.StatusInternalServerError, fmt.Errorf("failed to decode response: %w", err)
		}
	}

	return res.StatusCode, nil
}
