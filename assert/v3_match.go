package assert

import (
	"regexp"
	"strings"
)

// matchV3 compares actual output against a v3 pattern. nil means success.
// When normalizeNewlines is true, \r\n is rewritten to \n before matching.
// Same-name placeholder bindings must capture equal values across the template.
func matchV3(p v3Pattern, actual string, normalizeNewlines bool) error {
	if normalizeNewlines {
		actual = strings.ReplaceAll(actual, "\r\n", "\n")
	}
	if err := v3CheckTrailingNewline(p.trailingNewline, actual); err != nil {
		return err
	}
	lines := v3SplitActualLines(actual)
	bindings := make(map[string]string)
	cursor := 0
	for _, it := range p.items {
		n, err := v3MatchItem(it, lines, cursor, 0, bindings)
		if err != nil {
			return err
		}
		cursor = n
	}
	if cursor < len(lines) {
		return matchErr("unexpected extra output at line %d", cursor+1)
	}
	return nil
}

func v3CheckTrailingNewline(templateTrailing bool, actual string) error {
	actualTrailing := len(actual) > 0 && actual[len(actual)-1] == '\n'
	if templateTrailing != actualTrailing {
		return matchErr("output mismatch: trailing newline policy violated")
	}
	return nil
}

func v3SplitActualLines(actual string) []string {
	if actual == "" {
		return nil
	}
	body := strings.TrimSuffix(actual, "\n")
	if body == "" && strings.HasSuffix(actual, "\n") {
		return []string{""}
	}
	return strings.Split(body, "\n")
}

func v3MatchItem(it v3Item, lines []string, cursor, lineBase int, bindings map[string]string) (int, error) {
	switch x := it.(type) {
	case v3RegexLine:
		return v3MatchRegexLine(x.Pattern, lines, cursor, lineBase, bindings)
	case v3OmitLine:
		return v3MatchOmitLine(x.Count, lines, cursor, lineBase)
	default:
		return cursor, matchErr("internal error: unknown v3 item type")
	}
}

func v3MatchRegexLine(pattern string, lines []string, cursor, lineBase int, bindings map[string]string) (int, error) {
	if cursor >= len(lines) {
		return cursor, matchErr("output mismatch at line %d: regex expected a line", lineBase+cursor+1)
	}
	re, err := regexp.Compile(pattern)
	if err != nil {
		return cursor, matchErr("internal error: invalid regex: %v", err)
	}
	line := lines[cursor]
	m := re.FindStringSubmatch(line)
	if m == nil {
		return cursor, matchErr("output mismatch at line %d:\n  regex %q did not match %q", lineBase+cursor+1, pattern, line)
	}
	names := re.SubexpNames()
	for i, name := range names {
		if name == "" || i >= len(m) {
			continue
		}
		val := m[i]
		if prev, ok := bindings[name]; ok && prev != val {
			return cursor, matchErr("placeholder %s binding mismatch: previously %q, now %q", name, prev, val)
		}
		bindings[name] = val
	}
	return cursor + 1, nil
}

func v3MatchOmitLine(count int, lines []string, cursor, lineBase int) (int, error) {
	if cursor+count > len(lines) {
		return cursor, matchErr("output mismatch at line %d: omit expected %d lines, got %d", lineBase+cursor+1, count, len(lines)-cursor)
	}
	return cursor + count, nil
}
