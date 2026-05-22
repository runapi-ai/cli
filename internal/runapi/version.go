package runapi

import (
	"runtime/debug"
	"strings"
)

// Version resolves from (in order):
//  1. ldflags -X github.com/runapi-ai/cli/internal/runapi.Version=<tag>
//     (goreleaser release binaries — debug.ReadBuildInfo cannot see the tag).
//  2. debug.ReadBuildInfo().Main.Version
//     (go install github.com/runapi-ai/cli/cmd/runapi@<tag-or-commit>).
//  3. "dev"
//     (local go run / go build with no tag).
var Version = "dev"

func init() {
	if Version != "dev" {
		Version = strings.TrimPrefix(Version, "v")
		return
	}
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return
	}
	v := info.Main.Version
	if v == "" || v == "(devel)" {
		return
	}
	Version = strings.TrimPrefix(v, "v")
}
