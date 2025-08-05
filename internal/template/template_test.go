package template

import (
	"errors"
	"testing"
)

func TestUpdate(t *testing.T) {
	t.Parallel()

	t.Run("fails to parse", func(t *testing.T) {
		t.Parallel()

		s := "{{ 1 + 1 }}"
		err := Update(&s, nil)

		if !errors.Is(err, ErrInvalidTemplate) {
			t.Errorf("expected error %v, got %v", ErrInvalidTemplate, err)
		}
	})

	t.Run("fails to execute", func(t *testing.T) {
		t.Parallel()

		s := "{{ .Data }}"
		err := Update(&s, struct{}{})

		if !errors.Is(err, ErrFailTemplate) {
			t.Errorf("expected error %v, got %v", ErrFailTemplate, err)
		}
	})

	t.Run("updates correctly", func(t *testing.T) {
		t.Parallel()

		s := "{{ .Data }}"
		err := Update(&s, struct{ Data string }{Data: "done"})

		if !errors.Is(err, nil) {
			t.Errorf("expected no error got %v", err)
		}

		if s != "done" {
			t.Errorf("expected string %q, got %q", "done", s)
		}
	})
}
