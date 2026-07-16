package assert

import (
	"github.com/xhd2015/doctest/assert/legacy_v1"
	"github.com/xhd2015/doctest/assert/legacy_v2"
)

// PatternSummary returns a stable summary string for test introspection.
func PatternSummary(p Pattern) string {
	switch p.kind {
	case patternV2:
		return legacy_v2.PatternSummary(p.v2pat)
	case patternV3:
		return v3PatternSummary(p.v3pat)
	default:
		return legacy_v1.PatternSummary(p.legacy)
	}
}
