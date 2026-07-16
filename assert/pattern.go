package assert

import (
	"github.com/xhd2015/doctest/assert/legacy_v1"
	"github.com/xhd2015/doctest/assert/legacy_v2"
)

type patternKind int

const (
	patternV1 patternKind = iota
	patternV2
	patternV3
)

// Pattern is an immutable parsed output template (v1, v2, or v3).
type Pattern struct {
	kind   patternKind
	legacy legacy_v1.Pattern
	v2pat  legacy_v2.Pattern
	v3pat  v3Pattern
}
