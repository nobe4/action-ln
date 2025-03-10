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
)

var errInvalidYAML = errors.New("invalid YAML")

type RawConfig struct {
	Links []RawLink `yaml:"links"`
}

type Config struct {
	Links       []Link            `json:"links"         yaml:"links"`
	LinksByRepo map[string][]Link `json:"links_by_repo"`
}

func Parse(r io.Reader) (Config, error) {
	rawC := RawConfig{}

	if err := yaml.
		NewDecoder(r, yaml.Strict()).
		Decode(&rawC); err != nil {
		return Config{}, fmt.Errorf("%w: %w", errInvalidYAML, err)
	}

	c := Config{}

	for _, l := range rawC.Links {
		link, err := parseLink(l)
		if err != nil {
			return Config{}, err
		}

		c.Links = append(c.Links, link)
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
