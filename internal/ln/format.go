package ln

import (
	"fmt"
	"strings"
	"text/template"

	"github.com/nobe4/action-ln/internal/config"
	"github.com/nobe4/action-ln/internal/environment"
)

const (
	branchName   = "auto-action-ln"
	pullTitle    = "auto(ln): update links"
	bodyTemplate = `
This automated PR updates the following files:
{{ $b := "` + "`" + `" }}

| From | To  |
| ---  | --- |
{{ range .Links -}}
| {{ $b }}{{ .From }}{{ $b }} | {{ $b }}{{ .To }}{{ $b }} |
{{ end }}

---

| Quick links | [execution]({{ .Environment.ExecURL }}) | [configuration]({{.Environment.Server}}{{ .Config.Source.HTMLPath }}) | [action-ln](https://github.com/nobe4/action-ln) |
| --- | --- | --- | --- |
`
)

func pullRequestBody(l config.Links, c *config.Config, e environment.Environment) (string, error) {
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
		Config:      c,
		Environment: e,
	}

	if err := t.Execute(&out, data); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return out.String(), nil
}
