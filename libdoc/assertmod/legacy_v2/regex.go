package legacy_v2

import (
	"regexp"
	"strings"
)

var (
	placeholderNameRE = regexp.MustCompile(`__[A-Z][A-Z0-9_]*__`)
	omitLineRE        = regexp.MustCompile(`^\.\.\.\s*(\d+)\s*lines\s+omitted\s*\.\.\.$`)
	omitIntentRE      = regexp.MustCompile(`^\.\.\.\s*.+\s+lines\s+omitted\s*\.\.\.$`)
)

func isOmitLine(line string) (int, bool) {
	m := omitLineRE.FindStringSubmatch(strings.TrimSpace(line))
	if m == nil {
		return 0, false
	}
	var count int
	for _, ch := range m[1] {
		count = count*10 + int(ch-'0')
	}
	return count, true
}

func hasRegexIntent(line string) bool {
	masked := maskProtectedRegions(line)
	return scanRegexSignals(masked)
}

func maskProtectedRegions(line string) string {
	out := line
	for _, name := range placeholderNameRE.FindAllString(line, -1) {
		out = strings.ReplaceAll(out, name, strings.Repeat(" ", len(name)))
	}
	for {
		start := strings.Index(out, ansiColorOpenPrefix)
		if start < 0 {
			break
		}
		openEnd := strings.Index(out[start:], ">")
		if openEnd < 0 {
			break
		}
		openEnd += start
		closeStart := strings.Index(out[openEnd+1:], ansiColorCloseTag)
		if closeStart < 0 {
			break
		}
		closeStart += openEnd + 1
		span := out[start : closeStart+len(ansiColorCloseTag)]
		out = strings.Replace(out, span, strings.Repeat(" ", len(span)), 1)
	}
	return out
}

func scanRegexSignals(s string) bool {
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case '.':
			if i+1 < len(s) {
				switch s[i+1] {
				case '*', '+', '?':
					return true
				case '.':
					return true
				}
			}
		case '^':
			if i == 0 {
				return true
			}
		case '$':
			if i == len(s)-1 {
				return true
			}
		case '\\':
			if i+1 < len(s) {
				switch s[i+1] {
				case 'd', 'D', 'w', 'W', 's', 'S', 'b', 'B':
					return true
				}
			}
		case '[':
			if findBalancedBracket(s, i) {
				return true
			}
		case '|':
			if isAlternationSignal(s, i) {
				return true
			}
		case '{':
			if isQuantifierBrace(s, i) {
				return true
			}
		case '*', '+':
			if i+1 < len(s) && s[i+1] == '?' {
				return true
			}
		}
	}
	return false
}

func findBalancedBracket(s string, start int) bool {
	if start >= len(s) || s[start] != '[' {
		return false
	}
	depth := 0
	for i := start; i < len(s); i++ {
		if s[i] == '[' {
			depth++
		} else if s[i] == ']' {
			depth--
			if depth == 0 {
				return true
			}
		}
	}
	return false
}

func isAlternationSignal(s string, pipeIdx int) bool {
	left := strings.TrimSpace(s[:pipeIdx])
	right := strings.TrimSpace(s[pipeIdx+1:])
	return left != "" && right != ""
}

func isQuantifierBrace(s string, start int) bool {
	if start >= len(s) || s[start] != '{' {
		return false
	}
	end := strings.Index(s[start:], "}")
	if end < 0 {
		return false
	}
	inner := s[start+1 : start+end]
	if inner == "" {
		return false
	}
	for i, ch := range inner {
		if ch == ',' {
			continue
		}
		if ch < '0' || ch > '9' {
			if i == 0 {
				return false
			}
			return false
		}
	}
	return true
}

func placeholderSubpattern(typ string) (string, error) {
	switch typ {
	case "string":
		return `[^\n]*?`, nil
	case "number":
		return `-?\d+(?:\.\d+)?`, nil
	default:
		return "", parseErr(-1, "unknown placeholder type %q", typ)
	}
}
