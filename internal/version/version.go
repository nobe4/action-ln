package version

import (
	"fmt"
	"runtime"
	"runtime/debug"
)

const template = "%s (%s) built for %s/%s at %s\nhttps://%s/releases/tag/%s"

// Leaving them global to allow setting them via ldflags.
// E.g. go build cmd/<path to>/main.go -ldflags "-X <path>/internal/version.tag=0.1.0".
//
//nolint:gochecknoglobals // see above.
var (
	tag    = "UNSET_TAG"
	commit = "UNSET_COMMIT"
	date   = "UNSET_DATE"
	path   = "UNSET_PATH"
	os     = "UNSET_OS"
	arch   = "UNSET_ARCH"
)

func String() string {
	if tag == "UNSET_TAG" &&
		commit == "UNSET_COMMIT" &&
		date == "UNSET_DATE" &&
		path == "UNSET_PATH" {
		parseBuildInfo()
	}

	return fmt.Sprintf(template, tag, commit, os, arch, date, path, tag)
}

//nolint:revive // This can difficultly be simplified.
func parseBuildInfo() {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return
	}

	if info.Main.Path != "" {
		path = info.Main.Path
	}

	if info.Main.Version != "" {
		tag = info.Main.Version
	}

	os = runtime.GOOS
	arch = runtime.GOARCH

	for _, setting := range info.Settings {
		if setting.Key == "vcs.revision" && setting.Value != "" {
			commit = setting.Value[:8]
		}

		if setting.Key == "vcs.time" && setting.Value != "" {
			date = setting.Value
		}
	}
}
