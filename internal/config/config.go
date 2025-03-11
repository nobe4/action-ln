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
	"os"

	"github.com/goccy/go-yaml"

	"github.com/nobe4/action-ln/internal/github"
)

var (
	errInvalidYAML  = errors.New("invalid YAML")
	errInvalidLinks = errors.New("invalid links")
)

type RawConfig struct {
	Defaults map[string]any `yaml:"defaults"`
	Links    []RawLink      `yaml:"links"`
}

type Config struct {
	Defaults Defaults `json:"defaults" yaml:"defaults"`
	Links    []Link   `json:"links"    yaml:"links"`
}

func New() *Config { return &Config{} }

func (c *Config) Parse(r io.Reader) error {
	rawC := RawConfig{}

	if err := yaml.
		NewDecoder(r, yaml.Strict()).
		Decode(&rawC); err != nil {
		return fmt.Errorf("%w: %w", errInvalidYAML, err)
	}

	var err error

	c.Defaults.parse(rawC.Defaults)

	if c.Links, err = c.parseLinks(rawC.Links); err != nil {
		return fmt.Errorf("%w: %w", errInvalidLinks, err)
	}

	return nil
}

func (c *Config) Populate(ctx context.Context, g github.FileGetter) error {
	for i, l := range c.Links {
		if err := l.populate(ctx, g); err != nil {
			return fmt.Errorf("failed to populate link %#v: %w", l, err)
		}

		c.Links[i] = l
	}

	return nil
}

func (c *Config) String() string {
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
