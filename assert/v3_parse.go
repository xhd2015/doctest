package assert

import (
	"fmt"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

var (
	v3PlaceholderNameRE = regexp.MustCompile(`__[A-Z][A-Z0-9_]*__`)
	v3OmitLineRE        = regexp.MustCompile(`^\.\.\.\s*(\d+)\s*lines\s+omitted\s*\.\.\.$`)
	v3OmitIntentRE      = regexp.MustCompile(`^\.\.\.\s*.+\s+lines\s+omitted\s*\.\.\.$`)
)

// dialectKind classifies a template for facade routing.
type dialectKind int

const (
	dialectNone dialectKind = iota // no YAML assert dialect → legacy_v1
	dialectV2
	dialectV3
	dialectUnknownVersion
)

// detectDialect classifies template routing:
//   - version: 2 → dialectV2
//   - version: 3 → dialectV3
//   - unknown version → dialectUnknownVersion (caller errors)
//   - YAML header with placeholder keys, no version → dialectV3 (default)
//   - otherwise → dialectNone (legacy_v1), including non-dialect --- frontmatter
func detectDialect(template string) (kind dialectKind, headerYAML, body string, versionStr string) {
	template = v3TrimLeadingBlankLinesBeforeHeader(template)
	if !strings.HasPrefix(template, "---") {
		return dialectNone, "", template, ""
	}
	rest := template[3:]
	if strings.HasPrefix(rest, "\r\n") {
		rest = rest[2:]
	} else if strings.HasPrefix(rest, "\n") {
		rest = rest[1:]
	} else {
		return dialectNone, "", template, ""
	}
	closeIdx := strings.Index(rest, "\n---")
	if closeIdx < 0 {
		return dialectNone, "", template, ""
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
		return dialectNone, "", template, ""
	}

	version, hasVersion := raw["version"]
	if hasVersion {
		versionStr = fmt.Sprint(version)
		// Normalize float-ish YAML ints (e.g. 2.0).
		switch v := version.(type) {
		case int:
			versionStr = fmt.Sprintf("%d", v)
		case int64:
			versionStr = fmt.Sprintf("%d", v)
		case float64:
			if v == float64(int(v)) {
				versionStr = fmt.Sprintf("%d", int(v))
			}
		}
		switch versionStr {
		case "2":
			return dialectV2, headerYAML, body, versionStr
		case "3":
			return dialectV3, headerYAML, body, versionStr
		default:
			return dialectUnknownVersion, headerYAML, body, versionStr
		}
	}

	// No version key: default to v3 only for YAML dialect with placeholder defs.
	// Plain frontmatter without placeholders stays on legacy_v1 (v2 missing-version leaf).
	for key := range raw {
		if v3PlaceholderNameRE.MatchString(key) {
			return dialectV3, headerYAML, body, ""
		}
	}
	return dialectNone, "", template, ""
}

// parseV3 parses a v3 YAML header and template body into a v3Pattern.
func parseV3(headerYAML, body string) (v3Pattern, error) {
	var raw map[string]any
	if err := yaml.Unmarshal([]byte(headerYAML), &raw); err != nil {
		return v3Pattern{}, parseErr(0, "invalid YAML header: %v", err)
	}

	placeholders := make(map[string]v3Placeholder)
	for key, val := range raw {
		if key == "version" {
			continue
		}
		if !v3PlaceholderNameRE.MatchString(key) {
			continue
		}
		text, ok := val.(string)
		if !ok {
			return v3Pattern{}, parseErr(0, "invalid placeholder definition for %s", key)
		}
		ph, err := v3ParsePlaceholderDef(key, text)
		if err != nil {
			return v3Pattern{}, err
		}
		placeholders[key] = ph
	}

	lines, trailing := v3SplitTemplateLines(body)
	var items []v3Item
	for i, line := range lines {
		pos := v3LineOffset(lines, i)
		if v3OmitIntentRE.MatchString(strings.TrimSpace(line)) {
			count, ok := v3IsOmitLine(line)
			if !ok {
				return v3Pattern{}, parseErr(pos, "invalid omit marker syntax")
			}
			items = append(items, v3OmitLine{Count: count})
			continue
		}
		it, err := v3ParseContentLine(line, pos, placeholders)
		if err != nil {
			return v3Pattern{}, err
		}
		items = append(items, it)
	}

	if err := v3ValidatePlaceholdersUsed(body, placeholders); err != nil {
		return v3Pattern{}, err
	}

	return v3Pattern{
		placeholders:    placeholders,
		items:           items,
		trailingNewline: trailing,
	}, nil
}

func v3ParsePlaceholderDef(name, raw string) (v3Placeholder, error) {
	parts := v3SplitCompactParts(raw)
	if len(parts) == 0 {
		return v3Placeholder{}, parseErr(0, "empty placeholder definition for %s", name)
	}

	meta := make(map[string]string)
	var typ string
	var customRE string
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
		switch key {
		case "type":
			typ = val
		case "regex":
			customRE = val
		default:
			meta[key] = val
		}
	}

	if typ != "" && customRE != "" {
		return v3Placeholder{}, parseErr(0, "placeholder %s cannot set both type= and regex=", name)
	}
	if typ == "" && customRE == "" {
		return v3Placeholder{}, parseErr(0, "placeholder %s missing type or regex", name)
	}

	shortName := strings.Trim(name, "_")

	if customRE != "" {
		if _, err := regexp.Compile(customRE); err != nil {
			return v3Placeholder{}, parseErr(0, "invalid regex for %s: %v", name, err)
		}
		return v3Placeholder{Name: shortName, Type: "regex", Regex: customRE, Metadata: meta}, nil
	}

	if typ != "string" && typ != "number" {
		return v3Placeholder{}, parseErr(0, "unknown placeholder type %q for %s", typ, name)
	}
	return v3Placeholder{Name: shortName, Type: typ, Metadata: meta}, nil
}

func v3SplitCompactParts(raw string) []string {
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

func v3ValidatePlaceholdersUsed(body string, defined map[string]v3Placeholder) error {
	used := v3PlaceholderNameRE.FindAllString(body, -1)
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

func v3IsOmitLine(line string) (int, bool) {
	m := v3OmitLineRE.FindStringSubmatch(strings.TrimSpace(line))
	if m == nil {
		return 0, false
	}
	var count int
	for _, ch := range m[1] {
		count = count*10 + int(ch-'0')
	}
	return count, true
}

// v3ParseContentLine builds a full-line raw RE with named capture groups for placeholders.
// Color spans keep CSI structure with QuoteMeta on open/inner/reset; outside tags stay raw RE.
func v3ParseContentLine(line string, pos int, placeholders map[string]v3Placeholder) (v3Item, error) {
	body, err := v3BuildRegexBody(line, pos, placeholders)
	if err != nil {
		return nil, err
	}
	pattern := "^" + body + "$"
	if _, err := regexp.Compile(pattern); err != nil {
		return nil, parseErr(pos, "invalid regex: %v", err)
	}
	return v3RegexLine{Pattern: pattern}, nil
}

func v3BuildRegexBody(line string, pos int, placeholders map[string]v3Placeholder) (string, error) {
	var b strings.Builder
	// Track which placeholder names already introduced a named group on this line
	// so repeats use backreferences (?P=NAME) — Go forbids duplicate named groups.
	namedIntroduced := make(map[string]bool)
	i := 0
	for i < len(line) {
		if strings.HasPrefix(line[i:], v3AnsiColorOpenPrefix) {
			openEnd := strings.Index(line[i:], ">")
			if openEnd < 0 {
				return "", parseErr(pos+i, "unclosed <ansi-color>")
			}
			openEnd += i
			openTag := line[i+1 : openEnd]
			if !strings.HasPrefix(openTag, v3AnsiColorTagName) {
				return "", parseErr(pos+i, "expected <ansi-color>")
			}
			specText := strings.TrimSpace(openTag[len(v3AnsiColorTagName):])
			spec, err := v3ParseColorSpec(specText)
			if err != nil {
				return "", err
			}
			closeIdx := strings.Index(line[openEnd+1:], v3AnsiColorCloseTag)
			if closeIdx < 0 {
				return "", parseErr(pos+i, "unclosed <ansi-color>")
			}
			inner := line[openEnd+1 : openEnd+1+closeIdx]
			if inner == "" {
				return "", parseErr(pos+i, "empty ansi-color inner text")
			}
			open, err := v3ResolveColorOpen(spec)
			if err != nil {
				return "", err
			}
			// CSI structure kept; inner text QuoteMeta'd (literal dots etc.).
			b.WriteString(regexp.QuoteMeta(open))
			b.WriteString(regexp.QuoteMeta(inner))
			b.WriteString(regexp.QuoteMeta(v3ColorReset))
			i = openEnd + 1 + closeIdx + len(v3AnsiColorCloseTag)
			continue
		}

		nextPH := v3PlaceholderNameRE.FindStringIndex(line[i:])
		nextColor := strings.Index(line[i:], v3AnsiColorOpenPrefix)
		litEnd := len(line)
		if nextPH != nil {
			litEnd = i + nextPH[0]
		}
		if nextColor >= 0 && i+nextColor < litEnd {
			litEnd = i + nextColor
		}
		if litEnd > i {
			// Raw RE syntax outside placeholders and color tags.
			b.WriteString(line[i:litEnd])
			i = litEnd
			continue
		}
		if nextPH != nil {
			full := line[i+nextPH[0] : i+nextPH[1]]
			ph, ok := placeholders[full]
			if !ok {
				short := strings.Trim(full, "_")
				return "", parseErr(pos+i, "undefined placeholder __%s__", short)
			}
			sub, err := v3PlaceholderSubpattern(ph)
			if err != nil {
				return "", err
			}
			if namedIntroduced[ph.Name] {
				// Same-name binding within the line via named backreference.
				b.WriteString("(?P=")
				b.WriteString(ph.Name)
				b.WriteString(")")
			} else {
				b.WriteString("(?P<")
				b.WriteString(ph.Name)
				b.WriteString(">")
				b.WriteString(sub)
				b.WriteString(")")
				namedIntroduced[ph.Name] = true
			}
			i += nextPH[1]
			continue
		}
		break
	}
	return b.String(), nil
}

func v3PlaceholderSubpattern(ph v3Placeholder) (string, error) {
	switch ph.Type {
	case "string":
		return `[^\n]*?`, nil
	case "number":
		return `-?\d+(?:\.\d+)?`, nil
	case "regex":
		if ph.Regex == "" {
			return "", parseErr(-1, "placeholder %s missing regex fragment", ph.Name)
		}
		return ph.Regex, nil
	default:
		return "", parseErr(-1, "unknown placeholder type %q", ph.Type)
	}
}

func v3SplitTemplateLines(s string) ([]string, bool) {
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

func v3LineOffset(lines []string, lineIdx int) int {
	if lineIdx <= 0 {
		return 0
	}
	offset := 0
	for i := 0; i < lineIdx && i < len(lines); i++ {
		offset += len(lines[i]) + 1
	}
	return offset
}

func v3TrimLeadingBlankLinesBeforeHeader(s string) string {
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
