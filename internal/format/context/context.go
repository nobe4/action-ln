package context

import (
	"fmt"
	"html/template"
	"strings"

	"github.com/nobe4/action-ln/internal/config"
	"github.com/nobe4/action-ln/internal/environment"
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

func (f Formatter) Format(tmpl string, data any) (string, error) {
	t, err := template.New("").Parse(tmpl)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	out := strings.Builder{}
	d := struct {
		Data        any
		Config      *config.Config
		Environment environment.Environment
	}{
		Data:        data,
		Config:      f.config,
		Environment: f.environment,
	}

	if err := t.Execute(&out, d); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return out.String(), nil
}
