package legacy_v1

import (
	"fmt"
	"regexp"
	"strings"
)

// MatchConfig controls legacy v1 matching behavior.
type MatchConfig struct {
	Contains          bool
	NormalizeNewlines bool
}

// Match compares actual output against pattern. nil means success.
func Match(p Pattern, actual string, cfg MatchConfig) error {
	if cfg.NormalizeNewlines {
		actual = strings.ReplaceAll(actual, "\r\n", "\n")
	}

	if cfg.Contains {
		return matchContainsMode(p, actual)
	}
	return matchExactMode(p, actual)
}

func matchContainsMode(p Pattern, actual string) error {
	for i := 0; i <= len(actual); i++ {
		if i == 0 || (i > 0 && actual[i-1] == '\n') {
			sub := actual[i:]
			if err := matchPrefixAt(p, sub, 0, false); err == nil {
				return nil
			}
		}
	}
	return matchErr("output does not contain expected contiguous subregion")
}

func matchExactMode(p Pattern, actual string) error {
	allowGaps := patternHasContains(p.Items)
	if !patternContainsOnly(p.Items) {
		if err := checkTrailingNewline(p.TrailingNewline, actual); err != nil {
			return err
		}
	}
	return matchExactAt(p, actual, 0, allowGaps, true)
}

// patternContainsOnly reports whether the pattern's top-level items are all
// ContainsBlocks (at least one) with no LiteralLine/PatternLine/RegexLine.
// Substring-only assertions should not care how the output terminates.
func patternContainsOnly(items []Item) bool {
	hasContains := false
	for _, item := range items {
		switch item.(type) {
		case ContainsBlock:
			hasContains = true
		case LiteralLine, PatternLine, RegexLine:
			return false
		}
	}
	return hasContains
}

func matchPrefixAt(p Pattern, actual string, lineBase int, allowGaps bool) error {
	return matchExactAt(p, actual, lineBase, allowGaps, false)
}

func matchExactAt(p Pattern, actual string, lineBase int, allowGaps bool, requireEOF bool) error {
	lines := splitActualLines(actual)
	st := &matchState{
		lines:     lines,
		lineBase:  lineBase,
		allowGaps: allowGaps,
	}

	for _, frag := range collectContainsFragments(p.Items) {
		if err := st.checkContainsFragment(frag); err != nil {
			return err
		}
	}

	cursor := 0
	for _, item := range p.Items {
		switch it := item.(type) {
		case ContainsBlock:
			continue
		default:
			n, err := st.matchItem(it, cursor)
			if err != nil {
				return err
			}
			cursor = n
		}
	}

	if requireEOF && !allowGaps && cursor < len(lines) {
		return matchErr("unexpected extra output at line %d", lineBase+cursor+1)
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

type matchState struct {
	lines     []string
	lineBase  int
	allowGaps bool
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

func patternHasContains(items []Item) bool {
	for _, item := range items {
		switch it := item.(type) {
		case ContainsBlock:
			return true
		case BlockOptional:
			if patternHasContains(it.Items) {
				return true
			}
		case AnyOfBlock:
			for _, b := range it.Branches {
				if patternHasContains(b.Items) {
					return true
				}
			}
		}
	}
	return false
}

func collectContainsFragments(items []Item) []ContainsFragment {
	var out []ContainsFragment
	for _, item := range items {
		switch it := item.(type) {
		case ContainsBlock:
			out = append(out, it.Fragments...)
		case BlockOptional:
			out = append(out, collectContainsFragments(it.Items)...)
		case AnyOfBlock:
			for _, b := range it.Branches {
				out = append(out, collectContainsFragments(b.Items)...)
			}
		}
	}
	return out
}

func (st *matchState) checkContainsFragment(frag ContainsFragment) error {
	for _, line := range st.lines {
		switch frag.Mode {
		case ContainsSubstring:
			if len(frag.Segments) > 0 {
				if _, err := st.matchSegments(frag.Segments, line, st.lineBase+1); err == nil {
					return nil
				}
				continue
			}
			if line == frag.Text || strings.Contains(line, frag.Text) {
				return nil
			}
		case ContainsStartWith:
			if strings.HasPrefix(line, frag.Text) || strings.Contains(line, frag.Text) {
				return nil
			}
		case ContainsEndWith:
			if strings.HasSuffix(line, frag.Text) {
				return nil
			}
		}
	}
	if len(frag.Segments) > 0 {
		return matchErr("missing contains pattern fragment")
	}
	return matchErr("missing contains fragment %q", frag.Text)
}

func (st *matchState) matchItem(item Item, cursor int) (int, error) {
	switch it := item.(type) {
	case LiteralLine:
		return st.matchLiteralLine(it.Text, cursor)
	case PatternLine:
		return st.matchPatternLine(it, cursor)
	case RegexLine:
		return st.matchRegexLine(it.Pattern, cursor)
	case BlockOptional:
		return st.matchBlockOptional(it, cursor)
	case AnyOfBlock:
		return st.matchAnyOfBlock(it, cursor)
	default:
		return cursor, matchErr("internal error: unknown item type")
	}
}

func (st *matchState) matchLiteralLine(want string, cursor int) (int, error) {
	if st.allowGaps {
		for i := cursor; i < len(st.lines); i++ {
			if st.lines[i] == want {
				return i + 1, nil
			}
		}
		return cursor, matchErr("output mismatch at line %d:\n  want: %q\n  got:  (not found)", st.lineBase+cursor+1, want)
	}
	if cursor >= len(st.lines) {
		return cursor, matchErr("output mismatch at line %d:\n  want: %q\n  got:  (end of output)", st.lineBase+cursor+1, want)
	}
	got := st.lines[cursor]
	if got != want {
		return cursor, matchErr("output mismatch at line %d:\n  want: %q\n  got:  %q", st.lineBase+cursor+1, want, got)
	}
	return cursor + 1, nil
}

func (st *matchState) matchPatternLine(pl PatternLine, cursor int) (int, error) {
	if cursor >= len(st.lines) {
		return cursor, matchErr("output mismatch at line %d: unexpected end of output", st.lineBase+cursor+1)
	}
	line := st.lines[cursor]
	rest, err := st.matchSegments(pl.Segments, line, st.lineBase+cursor+1)
	if err != nil {
		return cursor, err
	}
	if rest != "" {
		return cursor, matchErr("output mismatch at line %d:\n  unparsed remainder: %q", st.lineBase+cursor+1, rest)
	}
	return cursor + 1, nil
}

func (st *matchState) matchRegexLine(pattern string, cursor int) (int, error) {
	if cursor >= len(st.lines) {
		return cursor, matchErr("output mismatch at line %d: regex expected a line", st.lineBase+cursor+1)
	}
	re := regexp.MustCompile(pattern)
	if !re.MatchString(st.lines[cursor]) {
		return cursor, matchErr("output mismatch at line %d:\n  regex %q did not match %q", st.lineBase+cursor+1, pattern, st.lines[cursor])
	}
	return cursor + 1, nil
}

func (st *matchState) matchBlockOptional(bo BlockOptional, cursor int) (int, error) {
	if len(bo.Items) == 0 {
		return cursor, nil
	}
	end, err := st.matchItems(bo.Items, cursor, false)
	if err == nil {
		return end, nil
	}
	return cursor, nil
}

func (st *matchState) matchItems(items []Item, cursor int, required bool) (int, error) {
	pos := cursor
	for _, item := range items {
		n, err := st.matchItem(item, pos)
		if err != nil {
			if required {
				return pos, err
			}
			return cursor, err
		}
		pos = n
	}
	return pos, nil
}

func (st *matchState) matchAnyOfBlock(ao AnyOfBlock, cursor int) (int, error) {
	var reports []string
	for i, branch := range ao.Branches {
		end, err := st.matchItems(branch.Items, cursor, true)
		if err == nil {
			return end, nil
		}
		reports = append(reports, fmt.Sprintf("  branch %d (<expect>):\n    %s", i+1, indentMatchErr(err)))
	}
	msg := fmt.Sprintf("output mismatch at line %d: <any-of> — no branch matched\n  actual: %q\n%s",
		st.lineBase+cursor+1, st.lineAt(cursor), strings.Join(reports, "\n"))
	return cursor, &matchError{msg: msg}
}

func (st *matchState) lineAt(cursor int) string {
	if cursor >= len(st.lines) {
		return ""
	}
	return st.lines[cursor]
}

func indentMatchErr(err error) string {
	s := err.Error()
	s = strings.ReplaceAll(s, "\n", "\n    ")
	return s
}

func (st *matchState) matchSegments(segs []Segment, line string, lineNo int) (string, error) {
	rest := line
	for _, seg := range segs {
		var err error
		rest, err = st.matchSegment(seg, rest, lineNo)
		if err != nil {
			return rest, err
		}
	}
	return rest, nil
}

func (st *matchState) matchSegment(seg Segment, line string, lineNo int) (string, error) {
	switch s := seg.(type) {
	case Literal:
		if !strings.HasPrefix(line, s.Text) {
			return line, matchErr("output mismatch at line %d:\n  want prefix: %q\n  got:         %q", lineNo, s.Text, line)
		}
		return line[len(s.Text):], nil
	case Hint:
		if !strings.HasPrefix(line, s.Text) {
			return line, matchErr("output mismatch at line %d:\n  hint:%s expected %q, got %q", lineNo, s.Label, s.Text, line)
		}
		return line[len(s.Text):], nil
	case InlineOptional:
		if strings.HasPrefix(line, s.Text) {
			return line[len(s.Text):], nil
		}
		return line, nil
	case InlineAnyOf:
		return st.matchInlineAnyOf(s, line, lineNo)
	case AnsiColor:
		return st.matchAnsiColor(s, line, lineNo)
	case InlineRegex:
		return st.matchInlineRegex(s, line, lineNo)
	default:
		return line, matchErr("internal error: unknown segment type")
	}
}

func (st *matchState) matchInlineAnyOf(iao InlineAnyOf, line string, lineNo int) (string, error) {
	var reports []string
	for i, branch := range iao.Branches {
		rest, err := st.matchSegments(branch.Segments, line, lineNo)
		if err == nil && rest != "" {
			// branch matched prefix but left remainder — only valid if remainder consumed later
		}
		if err == nil {
			return rest, nil
		}
		reports = append(reports, fmt.Sprintf("  branch %d: %s", i+1, err.Error()))
	}
	return line, matchErr("output mismatch at line %d: inline <any-of> — no branch matched\n%s", lineNo, strings.Join(reports, "\n"))
}

func (st *matchState) matchInlineRegex(ir InlineRegex, line string, lineNo int) (string, error) {
	re := regexp.MustCompile(ir.Pattern)
	loc := re.FindStringIndex(line)
	if loc == nil || loc[0] != 0 {
		return line, matchErr("output mismatch at line %d:\n  regex %q did not match prefix of %q", lineNo, ir.Pattern, line)
	}
	return line[loc[1]:], nil
}

func (st *matchState) matchAnsiColor(ac AnsiColor, line string, lineNo int) (string, error) {
	open, err := resolveAnsiOpen(ac.Spec)
	if err != nil {
		return line, err
	}
	want := open + ac.Text + ansiReset
	if !strings.HasPrefix(line, want) {
		if strings.HasPrefix(line, ac.Text) {
			return line, matchErr("output mismatch at line %d:\n  ansi-color expected styled %q, got plain text", lineNo, ac.Text)
		}
		return line, matchErr("output mismatch at line %d:\n  ansi-color expected %q prefix, got %q", lineNo, want, line)
	}
	return line[len(want):], nil
}
