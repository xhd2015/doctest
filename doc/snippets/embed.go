package snippets

import _ "embed"

//go:embed DOCTEST_SPEC.md
var DocTestSpec string

//go:embed DOCTEST_DESIGN_SPEC.md
var DocTestDesignSpec string
