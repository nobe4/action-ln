package format

import (
	"fmt"
	"strings"
	"text/template"

	"github.com/nobe4/action-ln/internal/config"
	"github.com/nobe4/action-ln/internal/environment"
)

const (
	HeadBranch   = "auto-action-ln"
	PullTitle    = "auto(ln): update links"
	bodyTemplate = `
{{/* This defines a backtick character to use in the markdown. */}}
{{- $b := "` + "`" + `" -}}
This automated PR updates the following files:

| From | To  |
| ---  | --- |
{{ range .Links -}}
| {{ $b }}{{ .From }}{{ $b }} | {{ $b }}{{ .To }}{{ $b }} |
{{ end }}

---

| Quick links | [execution]({{ .Environment.ExecURL }}) | [configuration]({{ .Environment.Server }}{{ .Config.Source.HTMLPath }}) | [action-ln](https://github.com/nobe4/action-ln) |
| --- | --- | --- | --- |
`
)

type Formatter struct {
	config      *config.Config
	environment environment.Environment
}

func New(c *config.Config, e environment.Environment) Formatter {
	return Formatter{
		config:      c,
		environment: e,
	}
}

func (f Formatter) PullBody(l config.Links) (string, error) {
	t, err := template.New("Pull Body").Parse(bodyTemplate)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	out := strings.Builder{}
	data := struct {
		Links       config.Links
		Config      *config.Config
		Environment environment.Environment
	}{
		Links:       l,
		Config:      f.config,
		Environment: f.environment,
	}

	if err := t.Execute(&out, data); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return out.String(), nil
}
