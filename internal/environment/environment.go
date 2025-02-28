/*
Package environment implements helpers to get inputs and environment from the
GitHub action's environment variables.
*/
package environment

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

var (
	errInvalidEnvironment = errors.New("error parsing context")
	errNoToken            = errors.New("github token not found")
)

type Environment struct {
	Token string
}

func Parse() (Environment, error) {
	fmt.Fprintln(os.Stdout, "Environment variables:")

	for _, env := range os.Environ() {
		parts := strings.Split(env, "=")
		fmt.Fprintln(os.Stdout, parts[0])
	}

	e := Environment{}

	var err error

	if e.Token, err = parseToken(); err != nil {
		return e, fmt.Errorf("%w: %w", errInvalidEnvironment, err)
	}

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
