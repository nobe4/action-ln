/*
Package environment implements helpers to get inputs and environment from the
GitHub action's environment variables.
Called `environment` to avoid conflict with the `context` package.

https://docs.github.com/en/actions/writing-workflows/choosing-what-your-workflow-does/accessing-contextual-information-about-workflow-runs
*/
package environment

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/nobe4/action-ln/internal/github"
)

var (
	errInvalidEnvironment = errors.New("error parsing environment")
	errNoToken            = errors.New("github token not found")
	errNoRepo             = errors.New("github repository not found")
	errInvalidRepo        = errors.New("github repository invalid: want owner/repo")
)

const (
	defaultEndpoint = "https://api.github.com"
	defaultConfig   = ".github/ln-config.yaml"
	redacted        = "[redacted]"
)

type App struct {
	ID         string `json:"app_id"`          // INPUT_APP_ID
	PrivateKey string `json:"app_private_key"` // INPUT_APP_PRIVATE_KEY
	InstallID  string `json:"app_install_id"`  // INPUT_APP_INSTALL_ID
}

type Environment struct {
	Token    string      `json:"token"`    // GITHUB_TOKEN / INPUT_TOKEN
	Repo     github.Repo `json:"repo"`     // GITHUB_REPOSITORY
	Endpoint string      `json:"endpoint"` // GITHUB_API_URL
	Config   string      `json:"config"`   // INPUT_CONFIG
	App      App         `json:"app"`
}

//nolint:revive // No, I don't want to leak secrets.
func (e Environment) String() string {
	e.Token = redacted
	e.App.ID = redacted
	e.App.PrivateKey = redacted
	e.App.InstallID = redacted

	out, err := json.MarshalIndent(e, "", "  ")
	if err != nil {
		return err.Error()
	}

	return string(out)
}

func debug() {
	fmt.Fprintln(os.Stdout, "Environment variables:")

	for _, env := range os.Environ() {
		parts := strings.Split(env, "=")
		fmt.Fprintln(os.Stdout, parts[0])
	}
}

func Parse() (Environment, error) {
	debug()

	e := Environment{}

	var err error

	if e.Token, err = parseToken(); err != nil {
		return e, fmt.Errorf("%w: %w", errInvalidEnvironment, err)
	}

	if e.Repo, err = parseRepo(); err != nil {
		return e, fmt.Errorf("%w: %w", errInvalidEnvironment, err)
	}

	e.Endpoint = parseEndpoint()
	e.Config = parseConfig()
	e.App = parseApp()

	return e, nil
}

func parseToken() (string, error) {
	if token := os.Getenv("INPUT_TOKEN"); token != "" {
		return token, nil
	}

	if token := os.Getenv("GITHUB_TOKEN"); token != "" {
		return token, nil
	}

	return "", errNoToken
}

func parseRepo() (github.Repo, error) {
	repo := github.Repo{}
	repoName := os.Getenv("GITHUB_REPOSITORY")

	if repoName == "" {
		return repo, errNoRepo
	}

	var found bool
	repo.Owner.Login, repo.Repo, found = strings.Cut(repoName, "/")

	if !found {
		return repo, errInvalidRepo
	}

	return repo, nil
}

func parseEndpoint() string {
	if endpoint := os.Getenv("GITHUB_API_URL"); endpoint != "" {
		return endpoint
	}

	return defaultEndpoint
}

func parseConfig() string {
	if config := os.Getenv("INPUT_CONFIG"); config != "" {
		return config
	}

	return defaultConfig
}

func parseApp() App {
	return App{
		ID:         os.Getenv("INPUT_APP_ID"),
		PrivateKey: os.Getenv("INPUT_APP_PRIVATE_KEY"),
		InstallID:  os.Getenv("INPUT_APP_INSTALL_ID"),
	}
}
