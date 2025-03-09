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

	"github.com/nobe4/action-ln/internal/github"
)

var errInvalidYAML = errors.New("invalid YAML")

type RawConfig struct {
	Links []RawLink `yaml:"links"`
}

type RawLink struct {
	From any `yaml:"from"`
	To   any `yaml:"to"`
}

type Config struct {
	Links []Link `json:"links" yaml:"links"`
}

type Link struct {
	From github.File `json:"from" yaml:"from"`
	To   github.File `json:"to"   yaml:"to"`
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

func parseLink(raw RawLink) (Link, error) {
	from, err := parseFile(raw.From)
	if err != nil {
		return Link{}, err
	}

	to, err := parseFile(raw.To)
	if err != nil {
		return Link{}, err
	}

	return Link{From: from, To: to}, nil
}
