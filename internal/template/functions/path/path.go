package path

import (
	"path/filepath"
)

func TrimN(p string, n int) string {
	if n == 0 {
		return p
	}

	parts := splitAll(p)

	if n > 0 {
		parts = parts[n:]
	} else {
		parts = parts[:len(parts)+n]
	}

	return filepath.Join(parts...)
}

func splitAll(p string) []string {
	dir, last := filepath.Split(p)

	if dir == "" {
		return []string{last}
	}

	dir = filepath.Clean(dir)

	return append(splitAll(dir), last)
}
