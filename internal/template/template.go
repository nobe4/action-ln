package template

import (
	"errors"
	"fmt"
	"strings"
	"text/template"

	"github.com/nobe4/action-ln/internal/template/functions/path"
)

var (
	ErrInvalidTemplate = errors.New("invalid template")
	ErrFailTemplate    = errors.New("failed to execute template")
)

// Update replaces the parameter s with its content executed as a template.
func Update(s *string, data any) error {
	funcMap := map[string]any{
		"pathTrimN": path.TrimN,
	}

	t, err := template.
		New("").
		Funcs(funcMap).
		Parse(*s)
	if err != nil {
		return fmt.Errorf("%w: %q: %w", ErrInvalidTemplate, *s, err)
	}

	buf := strings.Builder{}
	if err := t.Execute(&buf, data); err != nil {
		return fmt.Errorf("%w: %q: %w", ErrFailTemplate, *s, err)
	}

	*s = buf.String()

	return nil
}
