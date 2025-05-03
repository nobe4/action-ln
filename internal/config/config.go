/*
Package config provides a permissive way to parse the config file.

It uses partial YAML unmarshalling to allow for a larger set of possible
configurations.
*/
package config

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"github.com/goccy/go-yaml"

	"github.com/nobe4/action-ln/internal/github"
	"github.com/nobe4/action-ln/internal/log"
)

var (
	errInvalidYAML     = errors.New("invalid YAML")
	errInvalidLinks    = errors.New("invalid links")
	errInvalidDefaults = errors.New("invalid defaults")
)

type RawConfig struct {
	Defaults RawDefaults `yaml:"defaults"`
	Links    []RawLink   `yaml:"links"`
}

type Config struct {
	Source   github.File `json:"source"   yaml:"source"`
	Defaults Defaults    `json:"defaults" yaml:"defaults"`
	Links    Links       `json:"links"    yaml:"links"`
}

func New(source github.File, repo github.Repo) *Config {
	return &Config{
		Source: source,
		Defaults: Defaults{
			Link: &Link{
				From: github.File{Repo: repo},
				To:   github.File{Repo: repo},
			},
		},
	}
}

func (c *Config) Parse(r io.Reader) error {
	var err error

	rawC := RawConfig{}

	if err = yaml.
		NewDecoder(r, yaml.Strict()).
		Decode(&rawC); err != nil {
		return fmt.Errorf("%w: %w", errInvalidYAML, err)
	}

	if err := c.parseDefaults(rawC.Defaults); err != nil {
		return fmt.Errorf("%w: %w", errInvalidDefaults, err)
	}

	if c.Links, err = c.parseLinks(rawC.Links); err != nil {
		return fmt.Errorf("%w: %w", errInvalidLinks, err)
	}

	return nil
}

func (c *Config) Populate(ctx context.Context, g github.FileGetter) error {
	log.Group("Populate config")
	defer log.GroupEnd()

	for i, l := range c.Links {
		if err := l.populate(ctx, g); err != nil {
			return fmt.Errorf("failed to populate link %#v: %w", l, err)
		}

		// TODO: check if this is really necessary.
		c.Links[i] = l
	}

	return nil
}

func (c *Config) String() string {
	out, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		log.Warn("Error marshaling config", "err", err)

		return fmt.Sprintf("%#v", c)
	}

	return string(out)
}

// TODO: refactor into `config/file/file.go`.
func getMapKey(m map[string]any, k string) string {
	if v, ok := m[k]; ok {
		if vs, ok := v.(string); ok {
			return vs
		}

		log.Debug("Value is not a string", "key", k, "value", v)
	} else {
		log.Debug("Key not found", "key", k)
	}

	return ""
}
