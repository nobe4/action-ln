package config

import (
	"context"
	"fmt"

	"github.com/nobe4/action-ln/internal/format"
	"github.com/nobe4/action-ln/internal/github"
	"github.com/nobe4/action-ln/internal/log"
)

type Links []*Link

func (l *Links) Equal(other []*Link) bool {
	if len(*l) != len(other) {
		return false
	}

	for i, link := range *l {
		if !link.Equal(other[i]) {
			return false
		}
	}

	return true
}

func (c *Config) parseLinks(raw []RawLink) (Links, error) {
	links := Links{}

	for i, rl := range raw {
		l, err := c.parseLink(rl)
		if err != nil {
			log.Debug("Failed to parse link", "index", i, "raw", rl, "error")

			return nil, err
		}

		links = append(links, l...)
	}

	return links, nil
}

func (l *Links) Update(
	ctx context.Context,
	g github.FileGetterUpdater,
	f format.Formatter,
	head github.Branch,
) (bool, error) {
	updated := false

	for _, link := range *l {
		if needUpdate, err := link.NeedUpdate(ctx, g, head); err != nil {
			return updated, fmt.Errorf("failed to check if link %q needs update: %w", link, err)
		} else if !needUpdate {
			log.Info("Update not needed", "link", link)

			continue
		}

		if err := link.Update(ctx, g, f, head); err != nil {
			return updated, fmt.Errorf("failed to process link %q: %w", l, err)
		}

		updated = true
	}

	return updated, nil
}

type Groups map[string]Links

func (l *Links) Groups() Groups {
	g := make(Groups)

	for _, link := range *l {
		g[link.To.Repo.String()] = append(g[link.To.Repo.String()], link)
	}

	return g
}
