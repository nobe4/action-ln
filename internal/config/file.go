package config

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/nobe4/action-ln/internal/github"
	"github.com/nobe4/action-ln/internal/log"
)

var (
	ErrInvalidFileType   = errors.New("invalid file type")
	ErrInvalidFileFormat = errors.New("invalid file format")
)

func (c *Config) parseFile(rawFile any) ([]github.File, error) {
	log.Debug("Parse file", "raw", rawFile)

	switch v := rawFile.(type) {
	case nil:
		return []github.File{}, nil

	case map[string]any:
		return c.parseFileMap(v)

	case string:
		return c.parseFileString(v)

	default:
		return []github.File{}, fmt.Errorf("%w: %v (%T)", ErrInvalidFileType, rawFile, rawFile)
	}
}

func (c *Config) parseFileMap(rawFile map[string]any) ([]github.File, error) {
	f := github.File{}

	f.Repo = parseRepoString(
		getMapKey(rawFile, "owner"),
		getMapKey(rawFile, "repo"),
	)
	if f.Repo.Empty() {
		f.Repo = c.Defaults.Repo
	}

	f.Path = getMapKey(rawFile, "path")
	f.Ref = getMapKey(rawFile, "ref")

	return []github.File{f}, nil
}

//nolint:funlen // This function doesn't need to be simplified.
func (c *Config) parseFileString(s string) ([]github.File, error) {
	// 'https://github.com/owner/repo/blob/ref/path/to/file'
	if m := regexp.
		MustCompile(`^https://github.com/(?P<owner>[\w-]+)/(?P<repo>[\w-]+)/blob/(?P<ref>[\w-]+)/(?P<path>.+)$`).
		FindStringSubmatch(s); len(m) > 0 {
		return []github.File{
			{
				Repo: github.Repo{
					Owner: github.User{Login: m[1]},
					Repo:  m[2],
				},
				Ref:  m[3],
				Path: m[4],
			},
		}, nil
	}

	// 'owner/repo/blob/ref/path/to/file'
	if m := regexp.
		MustCompile(`^(?P<owner>[\w-]+)/(?P<repo>[\w-]+)/blob/(?P<ref>[\w-]+)/(?P<path>.+)$`).
		FindStringSubmatch(s); len(m) > 0 {
		return []github.File{
			{
				Repo: github.Repo{
					Owner: github.User{Login: m[1]},
					Repo:  m[2],
				},
				Ref:  m[3],
				Path: m[4],
			},
		}, nil
	}

	// 'owner/repo:path/to/file@ref'
	if m := regexp.
		MustCompile(`^(?P<owner>[\w-]+)/(?P<repo>[\w-]+):(?P<path>.+)@(?P<ref>[\w-]+)$`).
		FindStringSubmatch(s); len(m) > 0 {
		return []github.File{
			{
				Repo: github.Repo{
					Owner: github.User{Login: m[1]},
					Repo:  m[2],
				},
				Path: m[3],
				Ref:  m[4],
			},
		}, nil
	}

	// 'path/to/file@ref'
	if m := regexp.
		MustCompile(`^(?P<path>[^@]+)@(?P<ref>[\w-]+)$`).
		FindStringSubmatch(s); len(m) > 0 {
		return []github.File{
			{
				Path: m[1],
				Ref:  m[2],
				Repo: c.Defaults.Repo,
			},
		}, nil
	}

	// 'path/to/file'
	if m := regexp.
		MustCompile(`^(?P<path>.+)$`).
		FindStringSubmatch(s); len(m) > 0 {
		return []github.File{
			{
				Path: m[1],
				Repo: c.Defaults.Repo,
			},
		}, nil
	}

	return []github.File{}, fmt.Errorf("%w: '%v'", ErrInvalidFileFormat, s)
}
