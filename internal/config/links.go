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

func (c *Config) parseLink(raw RawLink) (Links, error) {
	log.Group(fmt.Sprintf("Parse link: %+v", raw))
	defer log.GroupEnd()

	froms, err := c.parseFile(raw.From)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", errInvalidFrom, err)
	}

	tos, err := c.parseFile(raw.To)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", errInvalidTo, err)
	}

	links := combineLinks(froms, tos)

	links.FillDefaults(c.Defaults)
	links.FillMissing()

	if err := links.ApplyTemplate(c); err != nil {
		return nil, err
	}

	links.Filter()

	return links, nil
}

func combineLinks(froms, tos []github.File) Links {
	if len(froms) == 0 {
		log.Warn("Found no `from`, make sure you reference one.")

		return Links{}
	}

	useFrom := len(tos) == 0

	links := Links{}

	for _, from := range froms {
		if useFrom {
			// TODO: inherit the `ref` as well?
			tos = []github.File{{Path: from.Path}}
		}

		for _, to := range tos {
			links = append(links, &Link{From: from, To: to})
		}
	}

	return links
}

func (l *Links) FillMissing() {
	for _, l := range *l {
		l.fillMissing()
	}
}

func (l *Links) FillDefaults(d Defaults) {
	for _, l := range *l {
		l.fillDefaults(d)
	}
}

func (l *Links) ApplyTemplate(c *Config) error {
	for _, l := range *l {
		if err := l.applyTemplate(c); err != nil {
			return err
		}
	}

	return nil
}

func (l *Links) Filter() {
	newL := Links{}

	for _, l := range *l {
		if l.From.Equal(l.To) {
			log.Warn("Found moot link, ignoring", "link", l)

			continue
		}

		newL = append(newL, l)
	}

	*l = newL
}

func (l *Links) Update(
	ctx context.Context,
	g github.GetterUpdater,
	f format.Formatter,
	head github.Branch,
) bool {
	updated := false

	for _, link := range *l {
		needUpdate, err := link.NeedUpdate(ctx, g, head)
		if err != nil {
			log.Error("failed to check if link needs update", "link", link, "error", err)
			link.Status = StatusFailedToCheck

			continue
		}

		if !needUpdate {
			log.Info("Update not needed", "link", link)
			link.Status = StatusUpdateNotNeeded

			continue
		}

		if err := link.Update(ctx, g, f, head); err != nil {
			log.Error("failed to update", "link", link, "error", err)
			link.Status = StatusFailedToUpdate

			continue
		}

		updated = true
		link.Status = StatusUpdated
	}

	return updated
}

type Groups map[string]Links

func (l *Links) Groups() Groups {
	g := make(Groups)

	for _, link := range *l {
		g[link.To.Repo.String()] = append(g[link.To.Repo.String()], link)
	}

	return g
}

func (g Groups) String() string {
	out := ""

	for n, l := range g {
		out += fmt.Sprintf("Group %q:\n", n)

		for _, link := range l {
			out += fmt.Sprintf("\t%s\n", link)
		}
	}

	return out
}
