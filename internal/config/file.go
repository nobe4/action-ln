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

func (c *Config) parseFile(rawFile any) ([]github.File, error) {
	switch v := rawFile.(type) {
	case nil:
		return []github.File{}, nil

	case []any:
		return c.parseSlice(v)

	case map[string]any:
		return c.parseMap(v)

	case string:
		return c.parseString(v)

	default:
		return []github.File{}, fmt.Errorf("%w: %v (%T)", ErrInvalidFileType, rawFile, rawFile)
	}
}

func (c *Config) parseSlice(rawFiles []any) ([]github.File, error) {
	files := []github.File{}

	for _, rf := range rawFiles {
		f, err := c.parseFile(rf)
		if err != nil {
			return []github.File{}, err
		}

		files = append(files, f...)
	}

	return files, nil
}

func (c *Config) parseMap(rawFile map[string]any) ([]github.File, error) {
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

// NOTE: only the github and blob path require `owner` and `repo` to match, all
// others concerned can omit either/both.
// E.g. `/:path`, `owner/:path`, `/repo:path`, ...
// This is less readable, but very useful for testsing.
//
//nolint:revive // This function doesn't need to be simplified.
func (c *Config) parseString(s string) ([]github.File, error) {
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
		MustCompile(`^(?P<owner>[\w-]*)/(?P<repo>[\w-]*):(?P<path>[^@]+)@(?P<ref>[\w-]+)$`).
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

	// 'owner/repo:path/to/file'
	if m := regexp.
		MustCompile(`^(?P<owner>[\w-]*)/(?P<repo>[\w-]*):(?P<path>[^@]+)$`).
		FindStringSubmatch(s); len(m) > 0 {
		return []github.File{
			{
				Repo: github.Repo{
					Owner: github.User{Login: m[1]},
					Repo:  m[2],
				},
				Path: m[3],
			},
		}, nil
	}

	// 'owner/repo:@ref'
	if m := regexp.
		MustCompile(`^(?P<owner>[\w-]*)/(?P<repo>[\w-]*):@(?P<ref>[\w-]+)$`).
		FindStringSubmatch(s); len(m) > 0 {
		return []github.File{
			{
				Repo: github.Repo{
					Owner: github.User{Login: m[1]},
					Repo:  m[2],
				},
				Ref: m[3],
			},
		}, nil
	}

	// 'owner/repo:'
	if m := regexp.
		MustCompile(`^(?P<owner>[\w-]*)/(?P<repo>[\w-]*):$`).
		FindStringSubmatch(s); len(m) > 0 {
		return []github.File{
			{
				Repo: github.Repo{
					Owner: github.User{Login: m[1]},
					Repo:  m[2],
				},
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
