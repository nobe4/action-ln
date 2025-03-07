package config

import (
	"errors"
	"fmt"
	"strings"
)

var errInvalidFileType = errors.New("invalid file type")

type File struct {
	Repo  string `yaml:"repo"`
	Owner string `yaml:"owner"`
	Path  string `yaml:"path"`
}

func (f File) Equal(o File) bool {
	return f.Repo == o.Repo && f.Owner == o.Owner && f.Path == o.Path
}

func parseFile(rawFile any) (File, error) {
	switch v := rawFile.(type) {
	case map[string]any:
		return parseFileMap(v)

	default:
		return File{}, fmt.Errorf("%w: %v (%T)", errInvalidFileType, rawFile, rawFile)
	}
}

func parseFileMap(rawFile map[string]any) (File, error) {
	f := File{}

	f.Repo = getMapKey(rawFile, "repo")
	f.Owner = getMapKey(rawFile, "owner")
	f.Path = getMapKey(rawFile, "path")

	if strings.Contains(f.Repo, "/") {
		parts := strings.Split(f.Repo, "/")
		if len(parts) != 2 { //nolint:all // TODO: log that the repo is badly formatted
			// don't do anything
		}

		f.Owner = parts[0]
		f.Repo = parts[1]
	}

	return f, nil
}

func getMapKey(m map[string]any, k string) string {
	if v, ok := m[k]; ok {
		if vs, ok := v.(string); ok {
			return vs
		} else { //nolint:all // TODO: log that the key is not a string
		}
	} else { //nolint:all // TODO: log that the key is not found
	}

	return ""
}
