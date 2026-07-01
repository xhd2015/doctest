package assert

import (
	"regexp"
	"strings"
	"unicode"
)

// Parse parses template text into a Pattern.
func Parse(template string) (Pattern, error) {
	p := &parser{src: template}
	lines, trailing := splitTemplateLines(p.src)
	p.lines = lines
	p.idx = 0
	items, err := p.parseItems()
	if err != nil {
		return Pattern{}, err
	}
	return Pattern{items: items, trailingNewline: trailing}, nil
}

// MustParse parses template text or panics.
func MustParse(template string) Pattern {
	p, err := Parse(template)
	if err != nil {
		panic(err)
	}
	return p
}

type parser struct {
	src   string
	lines []string
	idx   int
}

func (p *parser) parseItems() ([]Item, error) {
	var items []Item
	for p.idx < len(p.lines) {
		item, err := p.parseItem()
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}

func splitTemplateLines(s string) ([]string, bool) {
	trailing := strings.HasSuffix(s, "\n")
	if s == "" {
		return nil, trailing
	}
	body := strings.TrimSuffix(s, "\n")
	if body == "" {
		return nil, trailing
	}
	return strings.Split(body, "\n"), trailing
}

func (p *parser) lineOffset(lineIdx int) int {
	if lineIdx <= 0 {
		return 0
	}
	offset := 0
	for i := 0; i < lineIdx && i < len(p.lines); i++ {
		offset += len(p.lines[i]) + 1
	}
	return offset
}

func (p *parser) parseItem() (Item, error) {
	line := p.lines[p.idx]
	trimmed := strings.TrimSpace(line)

	switch trimmed {
	case "<optional>":
		return p.parseBlockOptional()
	case "<any-of>":
		return p.parseAnyOfBlock()
	case "<contains>":
		return p.parseContainsBlock()
	case "<regex>":
		return p.parseBlockRegex()
	}

	pl, err := p.parseLine(line, p.lineOffset(p.idx))
	if err != nil {
		return nil, err
	}
	p.idx++
	return pl, nil
}

func (p *parser) parseBlockOptional() (Item, error) {
	p.idx++ // skip <optional>
	var inner []Item
	for p.idx < len(p.lines) {
		if strings.TrimSpace(p.lines[p.idx]) == "</optional>" {
			p.idx++
			return BlockOptional{Items: inner}, nil
		}
		item, err := p.parseItem()
		if err != nil {
			return nil, err
		}
		inner = append(inner, item)
	}
	return nil, parseErr(p.lineOffset(p.idx), "unclosed <optional>")
}

func (p *parser) parseAnyOfBlock() (Item, error) {
	p.idx++ // skip <any-of>
	var branches []ExpectBranch
	for p.idx < len(p.lines) {
		if strings.TrimSpace(p.lines[p.idx]) == "</any-of>" {
			p.idx++
			if len(branches) == 0 {
				return nil, parseErr(p.lineOffset(p.idx), "any-of requires at least one branch")
			}
			return AnyOfBlock{Branches: branches}, nil
		}
		if strings.TrimSpace(p.lines[p.idx]) != "<expect>" {
			return nil, parseErr(p.lineOffset(p.idx), "expected <expect> inside <any-of>")
		}
		p.idx++
		var body []Item
		for p.idx < len(p.lines) {
			if strings.TrimSpace(p.lines[p.idx]) == "</expect>" {
				p.idx++
				branches = append(branches, ExpectBranch{Items: body})
				break
			}
			item, err := p.parseItem()
			if err != nil {
				return nil, err
			}
			body = append(body, item)
		}
	}
	return nil, parseErr(p.lineOffset(p.idx), "unclosed <any-of>")
}

func (p *parser) parseContainsBlock() (Item, error) {
	p.idx++ // skip <contains>
	var fragments []ContainsFragment
	for p.idx < len(p.lines) {
		line := p.lines[p.idx]
		trimmed := strings.TrimSpace(line)
		if trimmed == "</contains>" {
			p.idx++
			return ContainsBlock{Fragments: fragments}, nil
		}
		if trimmed == "<start-with>" {
			frag, err := p.parseBlockStartWith()
			if err != nil {
				return nil, err
			}
			fragments = append(fragments, frag)
			continue
		}
		if trimmed == "<end-with>" {
			frag, err := p.parseBlockEndWith()
			if err != nil {
				return nil, err
			}
			fragments = append(fragments, frag)
			continue
		}
		if strings.Contains(line, "<start-with>") {
			frag, consumed, err := p.parseInlineStartWith(line, p.lineOffset(p.idx))
			if err != nil {
				return nil, err
			}
			if !consumed {
				return nil, parseErr(p.lineOffset(p.idx), "invalid start-with fragment")
			}
			fragments = append(fragments, frag)
			p.idx++
			continue
		}
		if strings.Contains(line, "<end-with>") {
			frag, consumed, err := p.parseInlineEndWith(line, p.lineOffset(p.idx))
			if err != nil {
				return nil, err
			}
			if !consumed {
				return nil, parseErr(p.lineOffset(p.idx), "invalid end-with fragment")
			}
			fragments = append(fragments, frag)
			p.idx++
			continue
		}
		parsed, err := p.parseLine(line, p.lineOffset(p.idx))
		if err != nil {
			return nil, err
		}
		switch it := parsed.(type) {
		case LiteralLine:
			fragments = append(fragments, ContainsFragment{Mode: ContainsSubstring, Text: it.Text})
		case PatternLine:
			fragments = append(fragments, ContainsFragment{Mode: ContainsSubstring, Segments: it.Segments})
		default:
			return nil, parseErr(p.lineOffset(p.idx), "unsupported contains fragment")
		}
		p.idx++
	}
	return nil, parseErr(p.lineOffset(p.idx), "unclosed <contains>")
}

func (p *parser) parseBlockStartWith() (ContainsFragment, error) {
	p.idx++
	var textLines []string
	for p.idx < len(p.lines) {
		if strings.TrimSpace(p.lines[p.idx]) == "</start-with>" {
			p.idx++
			return ContainsFragment{Mode: ContainsStartWith, Text: strings.Join(textLines, "\n")}, nil
		}
		textLines = append(textLines, p.lines[p.idx])
		p.idx++
	}
	return ContainsFragment{}, parseErr(p.lineOffset(p.idx), "unclosed <start-with>")
}

func (p *parser) parseBlockEndWith() (ContainsFragment, error) {
	p.idx++
	var textLines []string
	for p.idx < len(p.lines) {
		if strings.TrimSpace(p.lines[p.idx]) == "</end-with>" {
			p.idx++
			return ContainsFragment{Mode: ContainsEndWith, Text: strings.Join(textLines, "\n")}, nil
		}
		textLines = append(textLines, p.lines[p.idx])
		p.idx++
	}
	return ContainsFragment{}, parseErr(p.lineOffset(p.idx), "unclosed <end-with>")
}

func (p *parser) parseInlineStartWith(line string, pos int) (ContainsFragment, bool, error) {
	text, rest, ok, err := parseInlineWrapper(line, "start-with", pos)
	if err != nil || !ok {
		return ContainsFragment{}, false, err
	}
	if strings.TrimSpace(rest) != "" {
		return ContainsFragment{}, false, parseErr(pos, "start-with must be alone on line in block form or consume entire line inline")
	}
	return ContainsFragment{Mode: ContainsStartWith, Text: text}, true, nil
}

func (p *parser) parseInlineEndWith(line string, pos int) (ContainsFragment, bool, error) {
	text, rest, ok, err := parseInlineWrapper(line, "end-with", pos)
	if err != nil || !ok {
		return ContainsFragment{}, false, err
	}
	if strings.TrimSpace(rest) != "" {
		return ContainsFragment{}, false, parseErr(pos, "end-with must be alone on line in block form or consume entire line inline")
	}
	return ContainsFragment{Mode: ContainsEndWith, Text: text}, true, nil
}

func (p *parser) parseBlockRegex() (Item, error) {
	p.idx++ // skip <regex>
	if p.idx >= len(p.lines) {
		return nil, parseErr(p.lineOffset(p.idx), "regex block requires pattern line")
	}
	pattern := p.lines[p.idx]
	if _, err := regexp.Compile(pattern); err != nil {
		return nil, parseErr(p.lineOffset(p.idx), "invalid regex: %v", err)
	}
	p.idx++
	if p.idx >= len(p.lines) || strings.TrimSpace(p.lines[p.idx]) != "</regex>" {
		return nil, parseErr(p.lineOffset(p.idx), "unclosed <regex>")
	}
	p.idx++
	return RegexLine{Pattern: pattern}, nil
}

func (p *parser) parseLine(line string, pos int) (Item, error) {
	if !strings.Contains(line, "<") && !strings.Contains(line, "\\<") {
		return LiteralLine{Text: line}, nil
	}
	segments, err := parseSegments(line, pos)
	if err != nil {
		return nil, err
	}
	if len(segments) == 1 {
		if lit, ok := segments[0].(Literal); ok {
			if strings.Contains(line, "\\<") || strings.Contains(lit.Text, "<") {
				return PatternLine{Segments: segments}, nil
			}
			return LiteralLine{Text: lit.Text}, nil
		}
	}
	return PatternLine{Segments: segments}, nil
}

func parseSegments(line string, pos int) ([]Segment, error) {
	var segments []Segment
	i := 0
	for i < len(line) {
		if line[i] != '<' {
			litStart := i
			for i < len(line) {
				if line[i] == '\\' && i+1 < len(line) && line[i+1] == '<' {
					i += 2
					continue
				}
				if line[i] == '<' {
					break
				}
				i++
			}
			text := strings.ReplaceAll(line[litStart:i], "\\<", "<")
			segments = append(segments, Literal{Text: text})
			continue
		}
		if i > 0 && line[i-1] == '\\' {
			return nil, parseErr(pos+i, "invalid escape")
		}
		seg, n, err := parseTagAt(line, i, pos+i)
		if err != nil {
			return nil, err
		}
		segments = append(segments, seg)
		i += n
	}
	return segments, nil
}

func parseTagAt(line string, start, pos int) (Segment, int, error) {
	if !strings.HasPrefix(line[start:], "<") {
		return nil, 0, parseErr(pos, "expected tag")
	}
	closeIdx := strings.Index(line[start:], ">")
	if closeIdx < 0 {
		return nil, 0, parseErr(pos, "unclosed tag")
	}
	openTag := line[start+1 : start+closeIdx]

	switch {
	case openTag == "optional":
		return parseInlineClose(line, start, "optional", pos, func(inner string) (Segment, error) {
			return InlineOptional{Text: inner}, nil
		})
	case openTag == "any-of":
		return parseInlineAnyOf(line, start, pos)
	case strings.HasPrefix(openTag, "hint:"):
		label := openTag[5:]
		if !isValidLabel(label) {
			return nil, 0, parseErr(pos, "invalid hint label %q", label)
		}
		return parseHintSegment(line, start, label, pos)
	case strings.HasPrefix(openTag, "ansi-color"):
		specText := strings.TrimSpace(openTag[len("ansi-color"):])
		spec, err := parseAnsiSpec(specText)
		if err != nil {
			return nil, 0, err
		}
		return parseInlineClose(line, start, "ansi-color", pos, func(inner string) (Segment, error) {
			if inner == "" {
				return nil, parseErr(pos, "empty ansi-color inner text")
			}
			return AnsiColor{Spec: spec, Text: inner}, nil
		})
	case openTag == "regex":
		return parseInlineClose(line, start, "regex", pos, func(inner string) (Segment, error) {
			if inner == "" {
				return nil, parseErr(pos, "empty regex pattern")
			}
			if _, err := regexp.Compile(inner); err != nil {
				return nil, parseErr(pos, "invalid regex: %v", err)
			}
			return InlineRegex{Pattern: inner}, nil
		})
	default:
		if isRegisteredBareTag(openTag) {
			return nil, 0, parseErr(pos, "unknown or misplaced tag <%s>", openTag)
		}
		return nil, 0, parseErr(pos, "unknown tag <%s>", openTag)
	}
}

func parseInlineAnyOf(line string, start, pos int) (Segment, int, error) {
	openLen := len("<any-of>")
	rest := line[start+openLen:]
	var branches []InlineExpectBranch
	for {
		rest = strings.TrimLeft(rest, " \t")
		if !strings.HasPrefix(rest, "<expect>") {
			return nil, 0, parseErr(pos, "inline any-of requires <expect> branches")
		}
		rest = rest[len("<expect>"):]
		closeIdx := strings.Index(rest, "</expect>")
		if closeIdx < 0 {
			return nil, 0, parseErr(pos, "unclosed <expect>")
		}
		inner := rest[:closeIdx]
		seg, err := parseSegments(inner, pos)
		if err != nil {
			return nil, 0, err
		}
		branches = append(branches, InlineExpectBranch{Segments: seg})
		rest = rest[closeIdx+len("</expect>"):]
		rest = strings.TrimLeft(rest, " \t")
		if strings.HasPrefix(rest, "</any-of>") {
			rest = rest[len("</any-of>"):]
			consumed := len(line) - len(rest) - start
			return InlineAnyOf{Branches: branches}, consumed, nil
		}
	}
}

func parseHintSegment(line string, start int, label string, pos int) (Segment, int, error) {
	openEnd := strings.Index(line[start:], ">")
	if openEnd < 0 {
		return nil, 0, parseErr(pos, "unclosed tag")
	}
	openEnd += start
	rest := line[openEnd+1:]
	closePrefix := "</hint:"
	closeIdx := strings.Index(rest, closePrefix)
	if closeIdx < 0 {
		return nil, 0, parseErr(pos, "unclosed <hint:%s>", label)
	}
	closeRest := rest[closeIdx+len(closePrefix):]
	closeEnd := strings.Index(closeRest, ">")
	if closeEnd < 0 {
		return nil, 0, parseErr(pos, "unclosed <hint:%s>", label)
	}
	closeLabel := closeRest[:closeEnd]
	if closeLabel != label {
		return nil, 0, parseErr(pos, "hint label mismatch: open label %q, close label %q", label, closeLabel)
	}
	inner := rest[:closeIdx]
	if inner == "" {
		return nil, 0, parseErr(pos, "empty hint text")
	}
	consumed := openEnd + 1 + closeIdx + len(closePrefix) + closeEnd + 1 - start
	return Hint{Label: label, Text: inner}, consumed, nil
}

func parseInlineClose(line string, start int, closeName string, pos int, mk func(string) (Segment, error)) (Segment, int, error) {
	openEnd := strings.Index(line[start:], ">")
	if openEnd < 0 {
		return nil, 0, parseErr(pos, "unclosed tag")
	}
	openEnd += start
	closeTag := "</" + closeName + ">"
	closeIdx := strings.Index(line[openEnd+1:], closeTag)
	if closeIdx < 0 {
		return nil, 0, parseErr(pos, "unclosed <%s>", closeName)
	}
	inner := line[openEnd+1 : openEnd+1+closeIdx]
	seg, err := mk(inner)
	if err != nil {
		return nil, 0, err
	}
	consumed := openEnd + 1 + closeIdx + len(closeTag) - start
	return seg, consumed, nil
}

func parseInlineWrapper(line, name string, pos int) (text string, rest string, ok bool, err error) {
	open := "<" + name + ">"
	close := "</" + name + ">"
	start := strings.Index(line, open)
	if start < 0 {
		return "", line, false, nil
	}
	end := strings.Index(line, close)
	if end < 0 {
		return "", "", false, parseErr(pos+start, "unclosed <%s>", name)
	}
	text = line[start+len(open) : end]
	rest = line[:start] + line[end+len(close):]
	return text, rest, true, nil
}

func isValidLabel(label string) bool {
	if label == "" {
		return false
	}
	for i, r := range label {
		if i == 0 {
			if !unicode.IsLetter(r) {
				return false
			}
			continue
		}
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '_' {
			return false
		}
	}
	return true
}

func isRegisteredBareTag(name string) bool {
	switch name {
	case "id", "path", "cached", "run", "pass", "cwd", "optional", "any-of", "expect", "contains", "regex", "start-with", "end-with", "ansi-color":
		return true
	default:
		if strings.HasPrefix(name, "hint:") {
			return true
		}
		return false
	}
}
