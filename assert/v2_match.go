package assert

import (
	"regexp"
	"strings"
)

func matchV2(p v2Pattern, actual string, normalizeNewlines bool) error {
	if normalizeNewlines {
		actual = strings.ReplaceAll(actual, "\r\n", "\n")
	}
	if err := checkTrailingNewlineV2(p.trailingNewline, actual); err != nil {
		return err
	}
	lines := splitActualLinesV2(actual)
	cursor := 0
	for _, item := range p.items {
		n, err := matchV2Item(item, p.placeholders, lines, cursor, 0)
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

func checkTrailingNewlineV2(templateTrailing bool, actual string) error {
	actualTrailing := len(actual) > 0 && actual[len(actual)-1] == '\n'
	if templateTrailing != actualTrailing {
		return matchErr("output mismatch: trailing newline policy violated")
	}
	return nil
}

func splitActualLinesV2(actual string) []string {
	if actual == "" {
		return nil
	}
	body := strings.TrimSuffix(actual, "\n")
	if body == "" && strings.HasSuffix(actual, "\n") {
		return []string{""}
	}
	return strings.Split(body, "\n")
}

func matchV2Item(item v2Item, placeholders map[string]v2Placeholder, lines []string, cursor, lineBase int) (int, error) {
	switch it := item.(type) {
	case v2LiteralLine:
		return matchV2LiteralLine(it.Text, lines, cursor, lineBase)
	case v2PatternLine:
		return matchV2PatternLine(it, placeholders, lines, cursor, lineBase)
	case v2RegexLine:
		return matchV2RegexLine(it.Pattern, lines, cursor, lineBase)
	case v2OmitLine:
		return matchV2OmitLine(it.Count, lines, cursor, lineBase)
	default:
		return cursor, matchErr("internal error: unknown v2 item type")
	}
}

func matchV2LiteralLine(want string, lines []string, cursor, lineBase int) (int, error) {
	if cursor >= len(lines) {
		return cursor, matchErr("output mismatch at line %d:\n  want: %q\n  got:  (end of output)", lineBase+cursor+1, want)
	}
	got := lines[cursor]
	if got != want {
		return cursor, matchErr("output mismatch at line %d:\n  want: %q\n  got:  %q", lineBase+cursor+1, want, got)
	}
	return cursor + 1, nil
}

func matchV2PatternLine(pl v2PatternLine, placeholders map[string]v2Placeholder, lines []string, cursor, lineBase int) (int, error) {
	if cursor >= len(lines) {
		return cursor, matchErr("output mismatch at line %d: unexpected end of output", lineBase+cursor+1)
	}
	line := lines[cursor]
	rest, err := matchV2Segments(pl.Segments, placeholders, line, lineBase+cursor+1)
	if err != nil {
		return cursor, err
	}
	if rest != "" {
		return cursor, matchErr("output mismatch at line %d:\n  unparsed remainder: %q", lineBase+cursor+1, rest)
	}
	return cursor + 1, nil
}

func matchV2RegexLine(pattern string, lines []string, cursor, lineBase int) (int, error) {
	if cursor >= len(lines) {
		return cursor, matchErr("output mismatch at line %d: regex expected a line", lineBase+cursor+1)
	}
	re := regexp.MustCompile(pattern)
	if !re.MatchString(lines[cursor]) {
		return cursor, matchErr("output mismatch at line %d:\n  regex %q did not match %q", lineBase+cursor+1, pattern, lines[cursor])
	}
	return cursor + 1, nil
}

func matchV2OmitLine(count int, lines []string, cursor, lineBase int) (int, error) {
	if cursor+count > len(lines) {
		return cursor, matchErr("output mismatch at line %d: omit expected %d lines, got %d", lineBase+cursor+1, count, len(lines)-cursor)
	}
	return cursor + count, nil
}

func matchV2Segments(segs []v2Segment, placeholders map[string]v2Placeholder, line string, lineNo int) (string, error) {
	rest := line
	for _, seg := range segs {
		var err error
		rest, err = matchV2Segment(seg, placeholders, rest, lineNo)
		if err != nil {
			return rest, err
		}
	}
	return rest, nil
}

func matchV2Segment(seg v2Segment, placeholders map[string]v2Placeholder, line string, lineNo int) (string, error) {
	switch s := seg.(type) {
	case v2Literal:
		if !strings.HasPrefix(line, s.Text) {
			return line, matchErr("output mismatch at line %d:\n  want prefix: %q\n  got:         %q", lineNo, s.Text, line)
		}
		return line[len(s.Text):], nil
	case v2PlaceholderRef:
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
	case v2Color:
		return matchV2Color(s, line, lineNo)
	default:
		return line, matchErr("internal error: unknown v2 segment type")
	}
}

func findPlaceholder(placeholders map[string]v2Placeholder, shortName string) (v2Placeholder, bool) {
	for key, ph := range placeholders {
		if ph.Name == shortName || key == "__"+shortName+"__" {
			return ph, true
		}
	}
	return v2Placeholder{}, false
}

func matchV2Color(c v2Color, line string, lineNo int) (string, error) {
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