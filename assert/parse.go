package assert

import (
	"github.com/xhd2015/doctest/assert/legacy_v1"
	"github.com/xhd2015/doctest/assert/legacy_v2"
)

// Parse parses template text into a Pattern.
// Routing:
//   - version: 2 → legacy_v2
//   - version: 3 or YAML dialect without version (placeholder defs) → v3
//   - unknown version → parse error
//   - otherwise → legacy_v1
func Parse(template string) (Pattern, error) {
	kind, headerYAML, body, versionStr := detectDialect(template)
	switch kind {
	case dialectV2:
		v2pat, err := legacy_v2.Parse(headerYAML, body)
		if err != nil {
			return Pattern{}, err
		}
		return Pattern{kind: patternV2, v2pat: v2pat}, nil
	case dialectV3:
		v3pat, err := parseV3(headerYAML, body)
		if err != nil {
			return Pattern{}, err
		}
		return Pattern{kind: patternV3, v3pat: v3pat}, nil
	case dialectUnknownVersion:
		return Pattern{}, parseErr(0, "unknown assert template version %s", versionStr)
	default:
		legacy, err := legacy_v1.Parse(template)
		if err != nil {
			return Pattern{}, err
		}
		return Pattern{kind: patternV1, legacy: legacy}, nil
	}
}

// MustParse parses template text or panics.
func MustParse(template string) Pattern {
	p, err := Parse(template)
	if err != nil {
		panic(err)
	}
	return p
}
