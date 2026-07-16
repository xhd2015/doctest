package legacy_v2

import (
	"regexp"
	"strings"
)

// Match compares actual output against a v2 pattern. nil means success.
// When normalizeNewlines is true, \r\n is rewritten to \n before matching.
func Match(p Pattern, actual string, normalizeNewlines bool) error {
	if normalizeNewlines {
		actual = strings.ReplaceAll(actual, "\r\n", "\n")
	}
	if err := checkTrailingNewline(p.trailingNewline, actual); err != nil {
		return err
	}
	lines := splitActualLines(actual)
	cursor := 0
	for _, it := range p.items {
		n, err := matchItem(it, p.placeholders, lines, cursor, 0)
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

func checkTrailingNewline(templateTrailing bool, actual string) error {
	actualTrailing := len(actual) > 0 && actual[len(actual)-1] == '\n'
	if templateTrailing != actualTrailing {
		return matchErr("output mismatch: trailing newline policy violated")
	}
	return nil
}

func splitActualLines(actual string) []string {
	if actual == "" {
		return nil
	}
	body := strings.TrimSuffix(actual, "\n")
	if body == "" && strings.HasSuffix(actual, "\n") {
		return []string{""}
	}
	return strings.Split(body, "\n")
}

func matchItem(it item, placeholders map[string]placeholder, lines []string, cursor, lineBase int) (int, error) {
	switch x := it.(type) {
	case literalLine:
		return matchLiteralLine(x.Text, lines, cursor, lineBase)
	case patternLine:
		return matchPatternLine(x, placeholders, lines, cursor, lineBase)
	case regexLine:
		return matchRegexLine(x.Pattern, lines, cursor, lineBase)
	case omitLine:
		return matchOmitLine(x.Count, lines, cursor, lineBase)
	default:
		return cursor, matchErr("internal error: unknown v2 item type")
	}
}

func matchLiteralLine(want string, lines []string, cursor, lineBase int) (int, error) {
	if cursor >= len(lines) {
		return cursor, matchErr("output mismatch at line %d:\n  want: %q\n  got:  (end of output)", lineBase+cursor+1, want)
	}
	got := lines[cursor]
	if got != want {
		return cursor, matchErr("output mismatch at line %d:\n  want: %q\n  got:  %q", lineBase+cursor+1, want, got)
	}
	return cursor + 1, nil
}

func matchPatternLine(pl patternLine, placeholders map[string]placeholder, lines []string, cursor, lineBase int) (int, error) {
	if cursor >= len(lines) {
		return cursor, matchErr("output mismatch at line %d: unexpected end of output", lineBase+cursor+1)
	}
	line := lines[cursor]
	// Match the whole line as one regex built from segments. Segment-by-segment
	// matching of type=string used `^[^\n]*?` alone, which always matches the
	// empty string (non-greedy with no trailing constraint). Anchoring with $
	// lets non-greedy string placeholders expand correctly mid-line and at EOL.
	reBody, err := buildPatternLineRegexFromSegments(pl.Segments, placeholders)
	if err != nil {
		return cursor, err
	}
	re, err := regexp.Compile("^" + reBody + "$")
	if err != nil {
		return cursor, matchErr("internal error: invalid pattern line regex: %v", err)
	}
	if !re.MatchString(line) {
		// Fall back to sequential diagnostics for a more precise error message.
		if _, segErr := matchSegments(pl.Segments, placeholders, line, lineBase+cursor+1); segErr != nil {
			return cursor, segErr
		}
		return cursor, matchErr("output mismatch at line %d:\n  pattern did not match %q", lineBase+cursor+1, line)
	}
	return cursor + 1, nil
}

// buildPatternLineRegexFromSegments joins segment patterns into a single-line RE.
// Literals are quoted; placeholders use placeholderSubpattern (string is non-greedy).
func buildPatternLineRegexFromSegments(segs []segment, placeholders map[string]placeholder) (string, error) {
	var b strings.Builder
	for _, seg := range segs {
		switch s := seg.(type) {
		case literal:
			b.WriteString(regexp.QuoteMeta(s.Text))
		case placeholderRef:
			ph, ok := findPlaceholder(placeholders, s.Name)
			if !ok {
				return "", matchErr("internal error: undefined placeholder %s", s.Name)
			}
			sub, err := placeholderSubpattern(ph.Type)
			if err != nil {
				return "", err
			}
			b.WriteString("(?:")
			b.WriteString(sub)
			b.WriteString(")")
		case color:
			open, err := resolveColorOpen(s.Spec)
			if err != nil {
				return "", err
			}
			b.WriteString(regexp.QuoteMeta(open + s.Text + colorReset))
		default:
			return "", matchErr("internal error: unknown v2 segment type")
		}
	}
	return b.String(), nil
}

func matchRegexLine(pattern string, lines []string, cursor, lineBase int) (int, error) {
	if cursor >= len(lines) {
		return cursor, matchErr("output mismatch at line %d: regex expected a line", lineBase+cursor+1)
	}
	re := regexp.MustCompile(pattern)
	if !re.MatchString(lines[cursor]) {
		return cursor, matchErr("output mismatch at line %d:\n  regex %q did not match %q", lineBase+cursor+1, pattern, lines[cursor])
	}
	return cursor + 1, nil
}

func matchOmitLine(count int, lines []string, cursor, lineBase int) (int, error) {
	if cursor+count > len(lines) {
		return cursor, matchErr("output mismatch at line %d: omit expected %d lines, got %d", lineBase+cursor+1, count, len(lines)-cursor)
	}
	return cursor + count, nil
}

func matchSegments(segs []segment, placeholders map[string]placeholder, line string, lineNo int) (string, error) {
	rest := line
	for _, seg := range segs {
		var err error
		rest, err = matchSegment(seg, placeholders, rest, lineNo)
		if err != nil {
			return rest, err
		}
	}
	return rest, nil
}

func matchSegment(seg segment, placeholders map[string]placeholder, line string, lineNo int) (string, error) {
	switch s := seg.(type) {
	case literal:
		if !strings.HasPrefix(line, s.Text) {
			return line, matchErr("output mismatch at line %d:\n  want prefix: %q\n  got:         %q", lineNo, s.Text, line)
		}
		return line[len(s.Text):], nil
	case placeholderRef:
		ph, ok := findPlaceholder(placeholders, s.Name)
		if !ok {
			return line, matchErr("internal error: undefined placeholder %s", s.Name)
		}
		sub, err := placeholderSubpattern(ph.Type)
		if err != nil {
			return line, err
		}
		re := regexp.MustCompile("^" + sub)
		loc := re.FindStringIndex(line)
		if loc == nil || loc[0] != 0 {
			return line, matchErr("output mismatch at line %d:\n  placeholder __%s__ (%s) did not match %q", lineNo, s.Name, ph.Type, line)
		}
		return line[loc[1]:], nil
	case color:
		return matchColor(s, line, lineNo)
	default:
		return line, matchErr("internal error: unknown v2 segment type")
	}
}

func findPlaceholder(placeholders map[string]placeholder, shortName string) (placeholder, bool) {
	for key, ph := range placeholders {
		if ph.Name == shortName || key == "__"+shortName+"__" {
			return ph, true
		}
	}
	return placeholder{}, false
}

func matchColor(c color, line string, lineNo int) (string, error) {
	open, err := resolveColorOpen(c.Spec)
	if err != nil {
		return line, err
	}
	want := open + c.Text + colorReset
	if !strings.HasPrefix(line, want) {
		if strings.HasPrefix(line, c.Text) {
			return line, matchErr("output mismatch at line %d:\n  ansi-color expected styled %q, got plain text", lineNo, c.Text)
		}
		return line, matchErr("output mismatch at line %d:\n  ansi-color expected %q prefix, got %q", lineNo, want, line)
	}
	return line[len(want):], nil
}
