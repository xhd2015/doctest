package legacy_v2

import (
	"fmt"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

// Detect reports whether template is a version-2 document and returns its
// YAML header and body when so. err is reserved for future hard failures;
// currently non-v2 input returns ok=false with a nil error.
func Detect(template string) (ok bool, headerYAML, body string, err error) {
	template = trimLeadingBlankLinesBeforeHeader(template)
	if !strings.HasPrefix(template, "---") {
		return false, "", template, nil
	}
	rest := template[3:]
	if strings.HasPrefix(rest, "\r\n") {
		rest = rest[2:]
	} else if strings.HasPrefix(rest, "\n") {
		rest = rest[1:]
	} else {
		return false, "", template, nil
	}
	closeIdx := strings.Index(rest, "\n---")
	if closeIdx < 0 {
		return false, "", template, nil
	}
	headerYAML = rest[:closeIdx]
	body = rest[closeIdx+4:]
	if strings.HasPrefix(body, "\r\n") {
		body = body[2:]
	} else if strings.HasPrefix(body, "\n") {
		body = body[1:]
	}

	var raw map[string]any
	if err := yaml.Unmarshal([]byte(headerYAML), &raw); err != nil {
		return false, "", template, nil
	}
	version, hasVersion := raw["version"]
	if !hasVersion {
		return false, "", template, nil
	}
	switch v := version.(type) {
	case int:
		if v != 2 {
			return false, "", template, nil
		}
	case int64:
		if v != 2 {
			return false, "", template, nil
		}
	case float64:
		if int(v) != 2 {
			return false, "", template, nil
		}
	default:
		if fmt.Sprint(v) != "2" {
			return false, "", template, nil
		}
	}
	return true, headerYAML, body, nil
}

// Parse parses a v2 YAML header and template body into a Pattern.
func Parse(headerYAML, body string) (Pattern, error) {
	var raw map[string]any
	if err := yaml.Unmarshal([]byte(headerYAML), &raw); err != nil {
		return Pattern{}, parseErr(0, "invalid YAML header: %v", err)
	}

	placeholders := make(map[string]placeholder)
	for key, val := range raw {
		if key == "version" {
			continue
		}
		if !placeholderNameRE.MatchString(key) {
			continue
		}
		text, ok := val.(string)
		if !ok {
			return Pattern{}, parseErr(0, "invalid placeholder definition for %s", key)
		}
		ph, err := parsePlaceholderDef(key, text)
		if err != nil {
			return Pattern{}, err
		}
		placeholders[key] = ph
	}

	lines, trailing := splitTemplateLines(body)
	var items []item
	for i, line := range lines {
		pos := lineOffset(lines, i)
		if omitIntentRE.MatchString(strings.TrimSpace(line)) {
			count, ok := isOmitLine(line)
			if !ok {
				return Pattern{}, parseErr(pos, "invalid omit marker syntax")
			}
			items = append(items, omitLine{Count: count})
			continue
		}
		it, err := parseLine(line, pos, placeholders, hasRegexIntent(line))
		if err != nil {
			return Pattern{}, err
		}
		items = append(items, it)
	}

	if err := validatePlaceholdersUsed(body, placeholders); err != nil {
		return Pattern{}, err
	}

	return Pattern{
		placeholders:    placeholders,
		items:           items,
		trailingNewline: trailing,
	}, nil
}

// MustParse parses header+body or panics.
func MustParse(headerYAML, body string) Pattern {
	p, err := Parse(headerYAML, body)
	if err != nil {
		panic(err)
	}
	return p
}

func parsePlaceholderDef(name, raw string) (placeholder, error) {
	parts := splitCompactParts(raw)
	if len(parts) == 0 {
		return placeholder{}, parseErr(0, "empty placeholder definition for %s", name)
	}

	meta := make(map[string]string)
	var typ string
	for _, part := range parts {
		eq := strings.Index(part, "=")
		if eq <= 0 {
			break
		}
		key := strings.TrimSpace(part[:eq])
		val := strings.TrimSpace(part[eq+1:])
		if key == "" || val == "" {
			break
		}
		if key == "type" {
			typ = val
		} else {
			meta[key] = val
		}
	}

	if typ == "" {
		return placeholder{}, parseErr(0, "placeholder %s missing type", name)
	}
	if typ != "string" && typ != "number" {
		return placeholder{}, parseErr(0, "unknown placeholder type %q for %s", typ, name)
	}

	shortName := strings.Trim(name, "_")
	return placeholder{Name: shortName, Type: typ, Metadata: meta}, nil
}

func splitCompactParts(raw string) []string {
	var parts []string
	var cur strings.Builder
	depth := 0
	for _, ch := range raw {
		switch ch {
		case '(':
			depth++
			cur.WriteRune(ch)
		case ')':
			if depth > 0 {
				depth--
			}
			cur.WriteRune(ch)
		case ',':
			if depth == 0 {
				parts = append(parts, strings.TrimSpace(cur.String()))
				cur.Reset()
				continue
			}
			cur.WriteRune(ch)
		default:
			cur.WriteRune(ch)
		}
	}
	if s := strings.TrimSpace(cur.String()); s != "" {
		parts = append(parts, s)
	}
	return parts
}

func validatePlaceholdersUsed(body string, defined map[string]placeholder) error {
	used := placeholderNameRE.FindAllString(body, -1)
	seen := make(map[string]bool)
	for _, name := range used {
		if seen[name] {
			continue
		}
		seen[name] = true
		if _, ok := defined[name]; !ok {
			short := strings.Trim(name, "_")
			return parseErr(0, "undefined placeholder __%s__", short)
		}
	}
	return nil
}

func parseLine(line string, pos int, placeholders map[string]placeholder, regexIntent bool) (item, error) {
	if regexIntent {
		pattern, err := buildRegexPattern(line, placeholders)
		if err != nil {
			return nil, err
		}
		if _, err := regexp.Compile(pattern); err != nil {
			return nil, parseErr(pos, "invalid regex: %v", err)
		}
		return regexLine{Pattern: pattern}, nil
	}

	if !strings.Contains(line, "__") && !strings.Contains(line, ansiColorOpenPrefix) {
		return literalLine{Text: line}, nil
	}

	segments, err := parseSegments(line, pos, placeholders)
	if err != nil {
		return nil, err
	}
	if len(segments) == 1 {
		if lit, ok := segments[0].(literal); ok {
			return literalLine{Text: lit.Text}, nil
		}
	}
	return patternLine{Segments: segments}, nil
}

func parseSegments(line string, pos int, placeholders map[string]placeholder) ([]segment, error) {
	var segments []segment
	i := 0
	for i < len(line) {
		if strings.HasPrefix(line[i:], ansiColorOpenPrefix) {
			seg, n, err := parseColorSegment(line, i, pos+i)
			if err != nil {
				return nil, err
			}
			segments = append(segments, seg)
			i += n
			continue
		}
		nextPH := placeholderNameRE.FindStringIndex(line[i:])
		nextColor := strings.Index(line[i:], ansiColorOpenPrefix)
		litEnd := len(line)
		if nextPH != nil {
			litEnd = i + nextPH[0]
		}
		if nextColor >= 0 && i+nextColor < litEnd {
			litEnd = i + nextColor
		}
		if litEnd > i {
			segments = append(segments, literal{Text: line[i:litEnd]})
			i = litEnd
			continue
		}
		if nextPH != nil {
			name := line[i+nextPH[0] : i+nextPH[1]]
			if _, ok := placeholders[name]; !ok {
				short := strings.Trim(name, "_")
				return nil, parseErr(pos+i, "undefined placeholder __%s__", short)
			}
			short := strings.Trim(name, "_")
			segments = append(segments, placeholderRef{Name: short})
			i += nextPH[1]
			continue
		}
		if nextColor >= 0 {
			continue
		}
		break
	}
	return segments, nil
}

func parseColorSegment(line string, start, pos int) (segment, int, error) {
	openEnd := strings.Index(line[start:], ">")
	if openEnd < 0 {
		return nil, 0, parseErr(pos, "unclosed <ansi-color>")
	}
	openEnd += start
	openTag := line[start+1 : openEnd]
	if !strings.HasPrefix(openTag, ansiColorTagName) {
		return nil, 0, parseErr(pos, "expected <ansi-color>")
	}
	specText := strings.TrimSpace(openTag[len(ansiColorTagName):])
	spec, err := parseColorSpec(specText)
	if err != nil {
		return nil, 0, err
	}
	closeIdx := strings.Index(line[openEnd+1:], ansiColorCloseTag)
	if closeIdx < 0 {
		return nil, 0, parseErr(pos, "unclosed <ansi-color>")
	}
	inner := line[openEnd+1 : openEnd+1+closeIdx]
	if inner == "" {
		return nil, 0, parseErr(pos, "empty ansi-color inner text")
	}
	consumed := openEnd + 1 + closeIdx + len(ansiColorCloseTag) - start
	return color{Spec: spec, Text: inner}, consumed, nil
}

func buildRegexPattern(line string, placeholders map[string]placeholder) (string, error) {
	var b strings.Builder
	i := 0
	for i < len(line) {
		if strings.HasPrefix(line[i:], ansiColorOpenPrefix) {
			openEnd := strings.Index(line[i:], ">")
			if openEnd < 0 {
				return "", parseErr(-1, "unclosed <ansi-color>")
			}
			openEnd += i
			openTag := line[i+1 : openEnd]
			specText := strings.TrimSpace(strings.TrimPrefix(openTag, ansiColorTagName))
			spec, err := parseColorSpec(specText)
			if err != nil {
				return "", err
			}
			closeIdx := strings.Index(line[openEnd+1:], ansiColorCloseTag)
			if closeIdx < 0 {
				return "", parseErr(-1, "unclosed <ansi-color>")
			}
			inner := line[openEnd+1 : openEnd+1+closeIdx]
			open, err := resolveColorOpen(spec)
			if err != nil {
				return "", err
			}
			b.WriteString(regexp.QuoteMeta(open))
			b.WriteString(regexp.QuoteMeta(inner))
			b.WriteString(regexp.QuoteMeta(colorReset))
			i = openEnd + 1 + closeIdx + len(ansiColorCloseTag)
			continue
		}
		nextPH := placeholderNameRE.FindStringIndex(line[i:])
		nextColor := strings.Index(line[i:], ansiColorOpenPrefix)
		litEnd := len(line)
		if nextPH != nil {
			litEnd = i + nextPH[0]
		}
		if nextColor >= 0 && i+nextColor < litEnd {
			litEnd = i + nextColor
		}
		if litEnd > i {
			// Regex-intent lines: non-placeholder text is raw RE syntax (e.g. .* or (a|b)).
			// Do not QuoteMeta here — that would turn intentional metacharacters into literals
			// and break pure-regex cookbook cases (V2-M4, V2-M13).
			b.WriteString(line[i:litEnd])
			i = litEnd
			continue
		}
		if nextPH != nil {
			name := line[i+nextPH[0] : i+nextPH[1]]
			ph, ok := placeholders[name]
			if !ok {
				short := strings.Trim(name, "_")
				return "", parseErr(-1, "undefined placeholder __%s__", short)
			}
			sub, err := placeholderSubpattern(ph.Type)
			if err != nil {
				return "", err
			}
			b.WriteString(sub)
			i += nextPH[1]
			continue
		}
		break
	}
	return b.String(), nil
}

func splitTemplateLines(s string) ([]string, bool) {
	trailing := strings.HasSuffix(s, "\n")
	if s == "" {
		return nil, trailing
	}
	body := strings.TrimSuffix(s, "\n")
	if body == "" && trailing {
		return []string{""}, trailing
	}
	return strings.Split(body, "\n"), trailing
}

func lineOffset(lines []string, lineIdx int) int {
	if lineIdx <= 0 {
		return 0
	}
	offset := 0
	for i := 0; i < lineIdx && i < len(lines); i++ {
		offset += len(lines[i]) + 1
	}
	return offset
}

func trimLeadingBlankLinesBeforeHeader(s string) string {
	for len(s) > 0 {
		if strings.HasPrefix(s, "---") {
			return s
		}
		s = strings.TrimLeft(s, " \t")
		if strings.HasPrefix(s, "---") {
			return s
		}
		if strings.HasPrefix(s, "\r\n") {
			s = s[2:]
			continue
		}
		if strings.HasPrefix(s, "\n") {
			s = s[1:]
			continue
		}
		return s
	}
	return s
}
