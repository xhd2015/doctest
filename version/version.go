package version

import (
	_ "embed"
	"strings"
)

//go:embed VERSION.txt
var versionBytes []byte

// Version returns the canonical doctest spec version from cmd/doctest/VERSION.txt.
func Version() string {
	return strings.TrimSpace(string(versionBytes))
}