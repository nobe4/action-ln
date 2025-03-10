package config

import (
	"errors"
	"fmt"
	"regexp"

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

	f.Repo = parseRepoString(
		getMapKey(rawFile, "owner"),
		getMapKey(rawFile, "repo"),
	)

	f.Path = getMapKey(rawFile, "path")
	f.Ref = getMapKey(rawFile, "ref")

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
			Path: m[3],
			Ref:  m[4],
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

	return github.File{}, fmt.Errorf("%w: '%v'", ErrInvalidFileFormat, s)
}
