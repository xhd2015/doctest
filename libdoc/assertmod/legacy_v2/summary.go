package legacy_v2

import (
	"fmt"
	"sort"
	"strings"
)

// PatternSummary returns a stable summary string for test introspection.
func PatternSummary(p Pattern) string {
	var parts []string
	if len(p.placeholders) > 0 {
		parts = append(parts, summarizePlaceholders(p.placeholders))
	}
	for _, it := range p.items {
		parts = append(parts, summarizeItem(it))
	}
	return combineSummaryParts(parts)
}

func summarizePlaceholders(placeholders map[string]placeholder) string {
	names := make([]string, 0, len(placeholders))
	for _, ph := range placeholders {
		names = append(names, ph.Name)
	}
	sort.Strings(names)
	var defs []string
	for _, name := range names {
		for _, ph := range placeholders {
			if ph.Name == name {
				defs = append(defs, summarizePlaceholder(ph))
				break
			}
		}
	}
	return "Placeholders:" + strings.Join(defs, ",")
}

func summarizePlaceholder(ph placeholder) string {
	s := fmt.Sprintf("%s{%s", ph.Name, ph.Type)
	keys := make([]string, 0, len(ph.Metadata))
	for k := range ph.Metadata {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		s += "," + k + "=" + ph.Metadata[k]
	}
	s += "}"
	return s
}

func summarizeItem(it item) string {
	switch x := it.(type) {
	case literalLine:
		return "LiteralLine"
	case patternLine:
		return "PatternLine segments:" + summarizeSegments(x.Segments)
	case regexLine:
		return "RegexLine"
	case omitLine:
		return fmt.Sprintf("OmitLine{%d}", x.Count)
	default:
		return "Unknown"
	}
}

func summarizeSegments(segs []segment) string {
	parts := make([]string, 0, len(segs))
	for _, seg := range segs {
		parts = append(parts, summarizeSegment(seg))
	}
	return strings.Join(parts, "+")
}

func summarizeSegment(seg segment) string {
	switch s := seg.(type) {
	case literal:
		return "Literal"
	case placeholderRef:
		return fmt.Sprintf("Placeholder{%s}", s.Name)
	case color:
		return "Color " + strings.Join(s.Spec.Tokens, " ")
	default:
		return "Unknown"
	}
}

func combineSummaryParts(parts []string) string {
	if len(parts) == 0 {
		return ""
	}
	itemParts := parts
	if strings.HasPrefix(parts[0], "Placeholders:") {
		if len(parts) == 1 {
			return parts[0]
		}
		itemParts = parts[1:]
	}
	if len(itemParts) >= 2 && itemParts[0] == "LiteralLine" {
		allLiteral := true
		for _, p := range itemParts {
			if p != "LiteralLine" {
				allLiteral = false
				break
			}
		}
		if allLiteral {
			prefix := ""
			if strings.HasPrefix(parts[0], "Placeholders:") {
				prefix = parts[0] + "+"
			}
			return prefix + fmt.Sprintf("LiteralLine×%d", len(itemParts))
		}
	}
	if strings.HasPrefix(parts[0], "Placeholders:") {
		return strings.Join(parts, "+")
	}
	return strings.Join(itemParts, "+")
}
