package noop

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
)

type Client struct {
	fallback http.Client
}

func New() Client {
	return Client{fallback: *http.DefaultClient}
}

func (c Client) Do(req *http.Request) (*http.Response, error) {
	if req.Method == http.MethodGet {
		//nolint:wrapcheck // Get requests can be transparent.
		return c.fallback.Do(req)
	}

	return c.noop(req)
}

func (c Client) noop(req *http.Request) (*http.Response, error) {
	switch {
	// github.Auth
	case req.Method == http.MethodPost &&
		regexp.MustCompile("/app/installations/[^/]+/access_token").MatchString(req.URL.Path):
		//nolint:wrapcheck // Authentication must be transparent.
		return c.fallback.Do(req)

	// github.CreateBranch
	case req.Method == http.MethodPost &&
		regexp.MustCompile("/repos/[^/]+/[^/]+/git/refs").MatchString(req.URL.Path):
		return response(http.StatusCreated, `{}`), nil

	// github.UpdateFile
	case req.Method == http.MethodPut &&
		regexp.MustCompile("/repos/[^/]+/[^/]+/contents/.+").MatchString(req.URL.Path):
		return response(http.StatusOK, `{"sha":"noop_sha_1234"}`), nil

	// github.CreatePull
	case req.Method == http.MethodPost &&
		regexp.MustCompile("/repos/[^/]+/[^/]+/pulls").MatchString(req.URL.Path):
		return response(http.StatusOK, `{"number": -1}`), nil

	default:
		return response(http.StatusBadRequest, ""),
			fmt.Errorf("%w: %s path %s", errors.ErrUnsupported, req.Method, req.URL.Path)
	}
}

func response(status int, body string) *http.Response {
	return &http.Response{
		StatusCode: status,
		Body:       io.NopCloser(strings.NewReader(body)),
	}
}
