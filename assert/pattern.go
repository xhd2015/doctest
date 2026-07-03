package assert

import "github.com/xhd2015/doctest/assert/legacy_v1"

// Pattern is an immutable parsed output template (v1 or v2).
type Pattern struct {
	v2     bool
	legacy legacy_v1.Pattern
	v2pat  v2Pattern
}