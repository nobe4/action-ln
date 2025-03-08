package github

import (
	"context"
	"errors"
	"fmt"
	"net/http"
)

type AppToken struct {
	Token string `json:"token"`
}

var errGetAppToken = errors.New("failed to get app token")

// https://docs.github.com/en/rest/apps/apps?apiVersion=2022-11-28#create-an-installation-access-token-for-an-app
func (g GitHub) GetAppToken(ctx context.Context, install string, jwt string) (string, error) {
	t := AppToken{}
	path := fmt.Sprintf("/app/installations/%s/access_tokens", install)

	//nolint:revive // This is a temporary assignment becaues the JWT is needed
	// to get the App Token. Since g is not a pointer, the value won't propagate
	// to future calls.
	g.Token = jwt

	if _, err := g.req(ctx, http.MethodPost, path, nil, &t); err != nil {
		return "", fmt.Errorf("%w: %w", errGetAppToken, err)
	}

	return t.Token, nil
}
