package assert

import (
	"github.com/xhd2015/doctest/assert/legacy_v1"
)

// Match compares actual output against pattern. nil means success.
func Match(p Pattern, actual string, opts ...Option) error {
	cfg := applyOptions(opts...)
	if p.v2 {
		if cfg.mode == matchContains {
			return matchErr("MatchContains is not supported for v2 templates")
		}
		return matchV2(p.v2pat, actual, cfg.normalizeNewlines)
	}
	return legacy_v1.Match(p.legacy, actual, legacy_v1.MatchConfig{
		Contains:          cfg.mode == matchContains,
		NormalizeNewlines: cfg.normalizeNewlines,
	})
}