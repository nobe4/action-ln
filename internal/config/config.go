/*
Package config provides a permissive way to parse the config file.

It uses partial YAML unmarshalling to allow for a larger set of possible
configurations.
*/
package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/goccy/go-yaml"

	"github.com/nobe4/action-ln/internal/environment"
)

var errInvalidYAML = errors.New("invalid YAML")

type RawConfig struct {
	Defaults map[string]any `yaml:"defaults"`
	Links    []RawLink      `yaml:"links"`
}

type Config struct {
	Defaults Defaults `json:"defaults" yaml:"defaults"`
	Links    []Link   `json:"links"    yaml:"links"`
}

func Parse(r io.Reader, e environment.Environment) (Config, error) {
	rawC := RawConfig{}

	if err := yaml.
		NewDecoder(r, yaml.Strict()).
		Decode(&rawC); err != nil {
		return Config{}, fmt.Errorf("%w: %w", errInvalidYAML, err)
	}

	c := Config{}

	var err error

	c.Defaults.Repo = e.Repo
	c.Defaults.parse(rawC.Defaults)

	c.Links, err = parseLinks(rawC.Links)
	if err != nil {
		// TODO: add error
		return Config{}, err
	}

	return c, nil
}

func (c Config) String() string {
	out, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		fmt.Fprintln(os.Stdout, "Error marshaling config:", err)

		return fmt.Sprintf("%#v", c)
	}

	return string(out)
}

func getMapKey(m map[string]any, k string) string {
	if v, ok := m[k]; ok {
		if vs, ok := v.(string); ok {
			return vs
		} else { //nolint:all // TODO: log that the key is not a string
		}
	} else { //nolint:all // TODO: log that the key is not found
	}

	return ""
}
