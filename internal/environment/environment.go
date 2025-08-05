/*
Package environment implements helpers to get inputs and environment from the
GitHub action's environment variables.
Called `environment` to avoid conflict with the `context` package.

https://docs.github.com/en/actions/writing-workflows/choosing-what-your-workflow-does/accessing-contextual-information-about-workflow-runs
*/
package environment

import (
	"fmt"
	"os"
	"strings"

	"github.com/nobe4/gh-ln/pkg/environment"
)

func Parse() (environment.Environment, error) {
	e := environment.Environment{}

	var err error
	if e.Token, err = parseToken(); err != nil {
		return e, fmt.Errorf("%w: %w", environment.ErrInvalidEnvironment, err)
	}

	if e.Repo, err = environment.ParseRepo(os.Getenv("GITHUB_REPOSITORY")); err != nil {
		return e, fmt.Errorf("%w: %w", environment.ErrInvalidEnvironment, err)
	}

	e.Noop = parseNoop()
	e.Endpoint = parseEndpoint()
	e.Server = parseServer()
	e.RunID = parseRunID()
	e.Config = parseConfig()
	e.App = parseApp()
	e.OnAction = parseOnAction()
	e.Debug = parseDebug()
	e.LocalConfig = parseLocalConfig()

	e.ExecURL = fmt.Sprintf("%s/%s/actions/runs/%s", e.Server, e.Repo, e.RunID)

	return e, nil
}

func parseNoop() bool {
	return truthy(os.Getenv("INPUT_NOOP"))
}

func parseToken() (string, error) {
	if token := os.Getenv("INPUT_TOKEN"); token != "" {
		return token, nil
	}

	if token := os.Getenv("GITHUB_TOKEN"); token != "" {
		return token, nil
	}

	return "", environment.ErrNoToken
}

func parseEndpoint() string {
	if endpoint := os.Getenv("GITHUB_API_URL"); endpoint != "" {
		return endpoint
	}

	return environment.DefaultEndpoint
}

func parseServer() string {
	if server := os.Getenv("GITHUB_SERVER_URL"); server != "" {
		return server
	}

	return environment.DefaultServer
}

func parseRunID() string {
	if runID := os.Getenv("GITHUB_RUN_ID"); runID != "" {
		return runID
	}

	return environment.DefaultRunID
}

func parseConfig() string {
	if config := os.Getenv("INPUT_CONFIG"); config != "" {
		return config
	}

	return environment.DefaultConfig
}

func parseApp() environment.App {
	return environment.App{
		ID:         os.Getenv("INPUT_APP_ID"),
		PrivateKey: os.Getenv("INPUT_APP_PRIVATE_KEY"),
		InstallID:  os.Getenv("INPUT_APP_INSTALL_ID"),
	}
}

func parseOnAction() bool {
	return os.Getenv("GITHUB_RUN_ID") != ""
}

func parseDebug() bool {
	return truthy(os.Getenv("RUNNER_DEBUG"))
}

func parseLocalConfig() string {
	return os.Getenv("INPUT_LOCAL_CONFIG")
}

func truthy(s string) bool {
	switch strings.ToLower(s) {
	case "1", "true", "yes":
		return true
	default:
		return false
	}
}
