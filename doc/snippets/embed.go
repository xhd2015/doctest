package snippets

import (
	_ "embed"
	"strings"

	"github.com/xhd2015/doctest/version"
)

//go:embed DOCTEST_SPEC.md
var DocTestSpec string

//go:embed DOCTEST_DESIGN_SPEC.md
var DocTestDesignSpec string

func ResolvedDocTestSpec() string {
	return strings.ReplaceAll(DocTestSpec, "__DOCTEST_VERSION__", version.Version())
}
