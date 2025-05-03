package config

import (
	"github.com/nobe4/action-ln/internal/log"
)

type RawDefaults struct {
	Link RawLink `yaml:"link"`
}

type Defaults struct {
	Link *Link `json:"link" yaml:"link"`
}

func (d *Defaults) Equal(o *Defaults) bool {
	return d.Link.Equal(o.Link)
}

func (c *Config) parseDefaults(raw RawDefaults) error {
	log.Debug("Parse defaults", "raw", raw)

	links, err := c.parseLink(raw.Link)
	if err != nil {
		return err
	}

	switch len(links) {
	case 0:
	case 1:
		c.Defaults.Link = links[0]
	default:
		log.Warn("Defaults has more than one link, using the first", "links", links)
		c.Defaults.Link = links[0]
	}

	return nil
}

func (c *Config) fillMissing(l *Link) {
	if c.Defaults.Link != nil {
		c.fillDefaults(l)
	}

	// Set To from From
	if l.To.Repo.Empty() {
		l.To.Repo = l.From.Repo
	}

	if l.To.Path == "" {
		l.To.Path = l.From.Path
	}
}

func (c *Config) fillDefaults(l *Link) {
	if l.From.Repo.Empty() {
		l.From.Repo = c.Defaults.Link.From.Repo
	}

	if l.From.Path == "" {
		l.From.Path = c.Defaults.Link.From.Path
	}

	if l.To.Repo.Empty() {
		l.To.Repo = c.Defaults.Link.To.Repo
	}

	if l.To.Path == "" {
		l.To.Path = c.Defaults.Link.To.Path
	}
}
