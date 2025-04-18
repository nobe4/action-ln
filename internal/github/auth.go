package github

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/nobe4/action-ln/internal/jwt"
	"github.com/nobe4/action-ln/internal/log"
)

type AppToken struct {
	Token string `json:"token"`
}

var (
	errGetJWT      = errors.New("failed to get JWT")
	errGetAppToken = errors.New("failed to get app token")
)

func (g *GitHub) Auth(ctx context.Context, token, appID, appPrivateKey, appInstallID string) error {
	log.Group("Authentication")
	defer log.GroupEnd()

	g.Token = token

	if appID != "" && appPrivateKey != "" && appInstallID != "" {
		log.Info("Using app authentication")

		var jwtToken string

		jwtToken, err := jwt.New(time.Now().Unix(), appID, appPrivateKey)
		if err != nil {
			log.Error("Failed to create a JWT", "err", err)

			return fmt.Errorf("%w: %w", errGetJWT, err)
		}

		if g.Token, err = g.GetAppToken(ctx, appInstallID, jwtToken); err != nil {
			log.Error("Failed to get app token", "err", err)

			return err
		}
	} else {
		log.Info("Using token authentication")
	}

	return nil
}

// https://docs.github.com/en/rest/apps/apps?apiVersion=2022-11-28#create-an-installation-access-token-for-an-app
func (g *GitHub) GetAppToken(ctx context.Context, install string, jwtToken string) (string, error) {
	t := AppToken{}
	path := fmt.Sprintf("/app/installations/%s/access_tokens", install)

	token := g.Token
	defer func() { g.Token = token }()

	g.Token = jwtToken

	if _, err := g.req(ctx, http.MethodPost, path, nil, &t); err != nil {
		return "", fmt.Errorf("%w: %w", errGetAppToken, err)
	}

	return t.Token, nil
}
