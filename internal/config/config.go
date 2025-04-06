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
	errInvalidYAML  = errors.New("invalid YAML")
	errInvalidLinks = errors.New("invalid links")
)

type RawConfig struct {
	Defaults map[string]any `yaml:"defaults"`
	Links    []RawLink      `yaml:"links"`
}

type Config struct {
	Defaults Defaults `json:"defaults" yaml:"defaults"`
	Links    Links    `json:"links"    yaml:"links"`
}

func New() *Config { return &Config{} }

func (c *Config) Parse(r io.Reader) error {
	rawC := RawConfig{}

	if err := yaml.
		NewDecoder(r, yaml.Strict()).
		Decode(&rawC); err != nil {
		return fmt.Errorf("%w: %w", errInvalidYAML, err)
	}

	log.Debug("Parse defaults", "raw", rawC.Defaults)

	c.Defaults.parse(rawC.Defaults)

	log.Debug("Parse links", "raw", rawC.Links)

	var err error
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
		log.Warn("Error marshaling config", "err", err)

		return fmt.Sprintf("%#v", c)
	}

	return string(out)
}

func getMapKey(m map[string]any, k string) string {
	if v, ok := m[k]; ok {
		if vs, ok := v.(string); ok {
			return vs
		}

		log.Warn("Value is not a string", "key", k, "value", v)
	} else {
		log.Warn("Value not found", "key", k)
	}

	return ""
}
