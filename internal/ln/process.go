package ln

import (
	"context"
	"fmt"

	"github.com/nobe4/action-ln/internal/config"
	"github.com/nobe4/action-ln/internal/format"
	"github.com/nobe4/action-ln/internal/github"
	"github.com/nobe4/action-ln/internal/log"
)

const (
	headName         = "auto-action-ln"
	pullTitle        = "auto(ln): update links"
	pullBodyTemplate = `
{{/* This defines a backtick character to use in the markdown. */}}
{{- $b := "` + "`" + `" -}}
This automated PR updates the following files:

| From | To  |
| ---  | --- |
{{ range .Data -}}
| [{{ $b }}{{ .From }}{{ $b }}]({{ .From.HTMLURL }}) | {{ $b }}{{ .To.Path }}{{ $b }} |
{{ end }}

---

| Quick links | [execution]({{ .Environment.ExecURL }}) | [configuration]({{ .Environment.Server }}{{ .Config.Source.HTMLPath }}) | [action-ln](https://github.com/nobe4/action-ln) |
| --- | --- | --- | --- |
`
)

func processGroups(ctx context.Context, g *github.GitHub, f format.Formatter, groups config.Groups) error {
	for _, l := range groups {
		if err := processLinks(ctx, g, f, l); err != nil {
			return err
		}
	}

	return nil
}

//nolint:revive // Will try to refactor that later
func processLinks(ctx context.Context, g *github.GitHub, f format.Formatter, l config.Links) error {
	toRepo := l[0].To.Repo

	log.Group("Processing links " + toRepo.String())
	defer log.GroupEnd()

	base, head, err := g.GetBaseAndHeadBranches(ctx, toRepo, headName)
	if err != nil {
		return fmt.Errorf("failed to prepare branches: %w", err)
	}

	log.Debug("Parsed branches", "head", head, "base", base)

	updated, err := l.Update(ctx, g, f, head)
	if err != nil {
		return fmt.Errorf("failed to update the links: %w", err)
	}

	if !updated && head.New {
		log.Info("No link was updated, cleaning up.", "repo", toRepo, "branch", head.Name)

		if err := g.DeleteBranch(ctx, toRepo, head.Name); err != nil {
			return fmt.Errorf("failed to delete un-updated branch: %w", err)
		}

		return nil
	}

	pullBody, err := f.Format(pullBodyTemplate, l)
	if err != nil {
		return fmt.Errorf("failed to create pull request body: %w", err)
	}

	log.Debug("Pull body", "body", pullBody)

	pull, err := g.GetOrCreatePull(ctx, toRepo, base.Name, head.Name, pullTitle, pullBody)
	if err != nil {
		return fmt.Errorf("failed to get pull request: %w", err)
	}

	log.Info("Result pull request", "pull", pull, "new", pull.New)

	return nil
}
