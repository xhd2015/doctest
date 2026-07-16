package assert

import (
	"github.com/xhd2015/doctest/assert/legacy_v1"
	"github.com/xhd2015/doctest/assert/legacy_v2"
)

// Match compares actual output against pattern. nil means success.
func Match(p Pattern, actual string, opts ...Option) error {
	cfg := applyOptions(opts...)
	switch p.kind {
	case patternV2:
		if cfg.mode == matchContains {
			return matchErr("MatchContains is not supported for v2 templates")
		}
		return legacy_v2.Match(p.v2pat, actual, cfg.normalizeNewlines)
	case patternV3:
		if cfg.mode == matchContains {
			return matchErr("MatchContains is not supported for v3 templates")
		}
		return matchV3(p.v3pat, actual, cfg.normalizeNewlines)
	default:
		return legacy_v1.Match(p.legacy, actual, legacy_v1.MatchConfig{
			Contains:          cfg.mode == matchContains,
			NormalizeNewlines: cfg.normalizeNewlines,
		})
	}
}
