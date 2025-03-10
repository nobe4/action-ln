package config

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/nobe4/action-ln/internal/github"
)

var (
	ErrInvalidFileType   = errors.New("invalid file type")
	ErrInvalidFileFormat = errors.New("invalid file format")
)

func parseFile(rawFile any) (github.File, error) {
	switch v := rawFile.(type) {
	case map[string]any:
		return parseFileMap(v)
	case string:
		return parseFileString(v)

	default:
		return github.File{}, fmt.Errorf("%w: %v (%T)", ErrInvalidFileType, rawFile, rawFile)
	}
}

func parseFileMap(rawFile map[string]any) (github.File, error) {
	f := github.File{}

	f.Repo = github.Repo{
		Owner: github.User{
			Login: getMapKey(rawFile, "owner"),
		},
		Repo: getMapKey(rawFile, "repo"),
	}
	f.Path = getMapKey(rawFile, "path")
	f.Ref = getMapKey(rawFile, "ref")

	if strings.Contains(f.Repo.Repo, "/") {
		parts := strings.Split(f.Repo.Repo, "/")
		if len(parts) != 2 { //nolint:all // TODO: log that the repo is badly formatted
			// don't do anything
		}

		f.Repo.Owner.Login = parts[0]
		f.Repo.Repo = parts[1]
	}

	return f, nil
}

func parseFileString(s string) (github.File, error) {
	// 'https://github.com/owner/repo/blob/ref/path/to/file'
	if m := regexp.
		MustCompile(`^https://github.com/(?P<owner>[\w-]+)/(?P<repo>[\w-]+)/blob/(?P<ref>[\w-]+)/(?P<path>.+)$`).
		FindStringSubmatch(s); len(m) > 0 {
		return github.File{
			Repo: github.Repo{
				Owner: github.User{Login: m[1]},
				Repo:  m[2],
			},
			Ref:  m[3],
			Path: m[4],
		}, nil
	}

	// 'owner/repo/blob/ref/path/to/file'
	if m := regexp.
		MustCompile(`^(?P<owner>[\w-]+)/(?P<repo>[\w-]+)/blob/(?P<ref>[\w-]+)/(?P<path>.+)$`).
		FindStringSubmatch(s); len(m) > 0 {
		return github.File{
			Repo: github.Repo{
				Owner: github.User{Login: m[1]},
				Repo:  m[2],
			},
			Ref:  m[3],
			Path: m[4],
		}, nil
	}

	// 'owner/repo:path/to/file@ref'
	if m := regexp.
		MustCompile(`^(?P<owner>[\w-]+)/(?P<repo>[\w-]+):(?P<path>.+)@(?P<ref>[\w-]+)$`).
		FindStringSubmatch(s); len(m) > 0 {
		return github.File{
			Repo: github.Repo{
				Owner: github.User{Login: m[1]},
				Repo:  m[2],
			},
			Ref:  m[3],
			Path: m[4],
		}, nil
	}

	// 'path/to/file@ref'
	if m := regexp.
		MustCompile(`^(?P<path>[^@]+)@(?P<ref>[\w-]+)$`).
		FindStringSubmatch(s); len(m) > 0 {
		return github.File{Path: m[1], Ref: m[2]}, nil
	}

	// 'path/to/file'
	if m := regexp.
		MustCompile(`^(?P<path>.+)$`).
		FindStringSubmatch(s); len(m) > 0 {
		return github.File{Path: m[1]}, nil
	}

	return github.File{}, fmt.Errorf("%w: %v", ErrInvalidFileFormat, s)
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
