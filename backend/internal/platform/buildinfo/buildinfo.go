// Package buildinfo reports the running build's version.
package buildinfo

import "runtime/debug"

// Version is set at build time: -ldflags "-X github.com/an4eetos/inner_circle/backend/internal/platform/buildinfo.Version=v1.2.3".
var Version = ""

const shortRevisionLen = 12

// Get returns Version if set, otherwise the VCS revision embedded by the Go toolchain, otherwise "dev".
func Get() string {
	if Version != "" {
		return Version
	}
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return "dev"
	}
	var rev string
	var dirty bool
	for _, s := range info.Settings {
		switch s.Key {
		case "vcs.revision":
			rev = s.Value
		case "vcs.modified":
			dirty = s.Value == "true"
		default:
		}
	}
	if rev == "" {
		return "dev"
	}
	if len(rev) > shortRevisionLen {
		rev = rev[:shortRevisionLen]
	}
	if dirty {
		rev += "-dirty"
	}
	return rev
}
