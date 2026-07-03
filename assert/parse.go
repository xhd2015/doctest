package assert

import "github.com/xhd2015/doctest/assert/legacy_v1"

// Parse parses template text into a Pattern.
func Parse(template string) (Pattern, error) {
	isV2, headerYAML, body, err := detectV2(template)
	if err != nil {
		return Pattern{}, err
	}
	if isV2 {
		v2pat, err := parseV2(headerYAML, body)
		if err != nil {
			return Pattern{}, err
		}
		return Pattern{v2: true, v2pat: v2pat}, nil
	}
	legacy, err := legacy_v1.Parse(template)
	if err != nil {
		return Pattern{}, err
	}
	return Pattern{v2: false, legacy: legacy}, nil
}

// MustParse parses template text or panics.
func MustParse(template string) Pattern {
	p, err := Parse(template)
	if err != nil {
		panic(err)
	}
	return p
}