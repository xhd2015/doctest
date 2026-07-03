package assert

import "github.com/xhd2015/doctest/assert/legacy_v1"

// PatternSummary returns a stable summary string for test introspection.
func PatternSummary(p Pattern) string {
	if p.v2 {
		return summaryV2(p.v2pat)
	}
	return legacy_v1.PatternSummary(p.legacy)
}