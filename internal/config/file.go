package config

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

var (
	ErrInvalidFileType   = errors.New("invalid file type")
	ErrInvalidFileFormat = errors.New("invalid file format")
)

type File struct {
	Repo  string `yaml:"repo"`
	Owner string `yaml:"owner"`
	Path  string `yaml:"path"`
	Ref   string `yaml:"ref"`
}

func (f File) Equal(o File) bool {
	return f.Repo == o.Repo && f.Owner == o.Owner && f.Path == o.Path
}

func parseFile(rawFile any) (File, error) {
	switch v := rawFile.(type) {
	case map[string]any:
		return parseFileMap(v)
	case string:
		return parseFileString(v)

	default:
		return File{}, fmt.Errorf("%w: %v (%T)", ErrInvalidFileType, rawFile, rawFile)
	}
}

func parseFileMap(rawFile map[string]any) (File, error) {
	f := File{}

	f.Repo = getMapKey(rawFile, "repo")
	f.Owner = getMapKey(rawFile, "owner")
	f.Path = getMapKey(rawFile, "path")
	f.Ref = getMapKey(rawFile, "ref")

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

func parseFileString(s string) (File, error) {
	// 'https://github.com/owner/repo/blob/ref/path/to/file'
	if m := regexp.
		MustCompile(`^https://github.com/(?P<owner>[\w-]+)/(?P<repo>[\w-]+)/blob/(?P<ref>[\w-]+)/(?P<path>.+)$`).
		FindStringSubmatch(s); len(m) > 0 {
		return File{Owner: m[1], Repo: m[2], Ref: m[3], Path: m[4]}, nil
	}

	// 'owner/repo/blob/ref/path/to/file'
	if m := regexp.
		MustCompile(`^(?P<owner>[\w-]+)/(?P<repo>[\w-]+)/blob/(?P<ref>[\w-]+)/(?P<path>.+)$`).
		FindStringSubmatch(s); len(m) > 0 {
		return File{Owner: m[1], Repo: m[2], Ref: m[3], Path: m[4]}, nil
	}

	// 'owner/repo:path/to/file@ref'
	if m := regexp.
		MustCompile(`^(?P<owner>[\w-]+)/(?P<repo>[\w-]+):(?P<path>.+)@(?P<ref>[\w-]+)$`).
		FindStringSubmatch(s); len(m) > 0 {
		return File{Owner: m[1], Repo: m[2], Path: m[3], Ref: m[4]}, nil
	}

	// 'path/to/file@ref'
	if m := regexp.
		MustCompile(`^(?P<path>[^@]+)@(?P<ref>[\w-]+)$`).
		FindStringSubmatch(s); len(m) > 0 {
		return File{Path: m[1], Ref: m[2]}, nil
	}

	// 'path/to/file'
	if m := regexp.
		MustCompile(`^(?P<path>.+)$`).
		FindStringSubmatch(s); len(m) > 0 {
		return File{Path: m[1]}, nil
	}

	return File{}, fmt.Errorf("%w: %v", ErrInvalidFileFormat, s)
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
