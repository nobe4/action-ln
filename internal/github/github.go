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

	"github.com/nobe4/action-ln/internal/client"
	"github.com/nobe4/action-ln/internal/log"
)

type FileGetter interface {
	GetFile(ctx context.Context, f *File) error
}

type FileUpdater interface {
	UpdateFile(ctx context.Context, f File, head string, msg string) (File, error)
}

type FileGetterUpdater interface {
	FileGetter
	FileUpdater
}

var (
	ErrRequestFailed  = errors.New("request failed")
	ErrMarshalRequest = errors.New("failed to marshal request")
)

const (
	PathUser = "/user"
)

type GitHub struct {
	client   client.Doer
	Token    string
	endpoint string
}

func New(c client.Doer, endpoint string) *GitHub {
	return &GitHub{
		client:   c,
		endpoint: endpoint,
	}
}

type User struct {
	Login string `json:"login"`
}

func (g *GitHub) req(ctx context.Context, method, path string, body io.Reader, out any) (int, error) {
	url := g.endpoint + path

	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		log.Debug("Request", "method", method, "url", url, "status", "failed to create", "err", err)

		return http.StatusInternalServerError, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("Authorization", "Bearer "+g.Token)

	res, err := g.client.Do(req)
	if err != nil {
		log.Debug("Request", "method", method, "url", url, "err", err, "status", res.StatusCode)

		return res.StatusCode, fmt.Errorf("%w: %w", ErrRequestFailed, err)
	}
	defer res.Body.Close()

	log.Debug("Request", "method", method, "url", url, "status", res.StatusCode)

	code2XX := res.StatusCode >= http.StatusOK && res.StatusCode < http.StatusMultipleChoices
	if !code2XX {
		return res.StatusCode, fmt.Errorf("%w (%s %s): %s", ErrRequestFailed, method, url, res.Status)
	}

	if out != nil {
		if err := json.NewDecoder(res.Body).Decode(out); err != nil {
			return http.StatusInternalServerError, fmt.Errorf("failed to decode response: %w", err)
		}
	}

	return res.StatusCode, nil
}
