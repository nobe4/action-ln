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
	"github.com/nobe4/action-ln/internal/log"
)

var (
	ErrInvalidEnvironment = errors.New("error parsing environment")
	ErrNoToken            = errors.New("github token not found")
	ErrNoRepo             = errors.New("github repository not found")
	ErrInvalidRepo        = errors.New("github repository invalid: want owner/repo")
)

const (
	defaultEndpoint = "https://api.github.com"
	defaultServer   = "https://github.com"
	defaultConfig   = ".github/ln-config.yaml"
	defaultRunID    = ""
	redacted        = "[redacted]"
	missing         = "[missing]"
)

type App struct {
	ID         string `json:"app_id"`          // INPUT_APP_ID
	PrivateKey string `json:"app_private_key"` // INPUT_APP_PRIVATE_KEY
	InstallID  string `json:"app_install_id"`  // INPUT_APP_INSTALL_ID
}

type Environment struct {
	Noop     bool        `json:"noop"`     // INPUT_NOOP
	Token    string      `json:"token"`    // GITHUB_TOKEN / INPUT_TOKEN
	Repo     github.Repo `json:"repo"`     // GITHUB_REPOSITORY
	Server   string      `json:"server"`   // GITHUB_SERVER_URL
	Endpoint string      `json:"endpoint"` // GITHUB_API_URL
	RunID    string      `json:"run_id"`   // GITHUB_RUN_ID
	Config   string      `json:"config"`   // INPUT_CONFIG
	App      App         `json:"app"`
	OnAction bool        `json:"on_action"`
	ExecURL  string      `json:"exec_url"`
	Debug    bool        `json:"debug"` // RUNNER_DEBUG
}

//nolint:revive // No, I don't want to leak secrets.
func (e Environment) String() string {
	e.Token = missingOrRedacted(e.Token)
	e.App.ID = missingOrRedacted(e.App.ID)
	e.App.PrivateKey = missingOrRedacted(e.App.PrivateKey)
	e.App.InstallID = missingOrRedacted(e.App.InstallID)

	out, err := json.MarshalIndent(e, "", "  ")
	if err != nil {
		return err.Error()
	}

	return string(out)
}

func (e Environment) PrintDebug() {
	log.Info("Environment", "parsed", e)

	log.Group("Environment keys")
	defer log.GroupEnd()

	for _, env := range os.Environ() {
		parts := strings.Split(env, "=")
		log.Debug(parts[0])
	}
}

func Parse() (Environment, error) {
	e := Environment{}

	var err error

	if e.Token, err = parseToken(); err != nil {
		return e, fmt.Errorf("%w: %w", ErrInvalidEnvironment, err)
	}

	if e.Repo, err = parseRepo(); err != nil {
		return e, fmt.Errorf("%w: %w", ErrInvalidEnvironment, err)
	}

	e.Noop = parseNoop()
	e.Endpoint = parseEndpoint()
	e.Server = parseServer()
	e.RunID = parseRunID()
	e.Config = parseConfig()
	e.App = parseApp()
	e.OnAction = parseOnAction()
	e.Debug = parseDebug()

	e.ExecURL = fmt.Sprintf("%s/%s/actions/runs/%s", e.Server, e.Repo, e.RunID)

	return e, nil
}

func parseNoop() bool {
	return os.Getenv("INPUT_NOOP") == "true"
}

func parseToken() (string, error) {
	if token := os.Getenv("INPUT_TOKEN"); token != "" {
		return token, nil
	}

	if token := os.Getenv("GITHUB_TOKEN"); token != "" {
		return token, nil
	}

	return "", ErrNoToken
}

func parseRepo() (github.Repo, error) {
	repo := github.Repo{}
	repoName := os.Getenv("GITHUB_REPOSITORY")

	if repoName == "" {
		return repo, ErrNoRepo
	}

	var found bool
	repo.Owner.Login, repo.Repo, found = strings.Cut(repoName, "/")

	if !found {
		return repo, ErrInvalidRepo
	}

	return repo, nil
}

func parseEndpoint() string {
	if endpoint := os.Getenv("GITHUB_API_URL"); endpoint != "" {
		return endpoint
	}

	return defaultEndpoint
}

func parseServer() string {
	if server := os.Getenv("GITHUB_SERVER_URL"); server != "" {
		return server
	}

	return defaultServer
}

func parseRunID() string {
	if runID := os.Getenv("GITHUB_RUN_ID"); runID != "" {
		return runID
	}

	return defaultRunID
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

func parseOnAction() bool {
	return os.Getenv("GITHUB_RUN_ID") != ""
}

func parseDebug() bool {
	return os.Getenv("RUNNER_DEBUG") == "1"
}

func missingOrRedacted(s string) string {
	if s == "" {
		return missing
	}

	return redacted
}
