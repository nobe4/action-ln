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
