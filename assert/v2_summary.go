package assert

import (
	"fmt"
	"sort"
	"strings"
)

func summaryV2(p v2Pattern) string {
	var parts []string
	if len(p.placeholders) > 0 {
		parts = append(parts, summarizeV2Placeholders(p.placeholders))
	}
	for _, item := range p.items {
		parts = append(parts, summarizeV2Item(item))
	}
	return combineV2SummaryParts(parts)
}

func summarizeV2Placeholders(placeholders map[string]v2Placeholder) string {
	names := make([]string, 0, len(placeholders))
	for _, ph := range placeholders {
		names = append(names, ph.Name)
	}
	sort.Strings(names)
	var defs []string
	for _, name := range names {
		for _, ph := range placeholders {
			if ph.Name == name {
				defs = append(defs, summarizeV2Placeholder(ph))
				break
			}
		}
	}
	return "Placeholders:" + strings.Join(defs, ",")
}

func summarizeV2Placeholder(ph v2Placeholder) string {
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

func summarizeV2Item(item v2Item) string {
	switch it := item.(type) {
	case v2LiteralLine:
		return "LiteralLine"
	case v2PatternLine:
		return "PatternLine segments:" + summarizeV2Segments(it.Segments)
	case v2RegexLine:
		return "RegexLine"
	case v2OmitLine:
		return fmt.Sprintf("OmitLine{%d}", it.Count)
	default:
		return "Unknown"
	}
}

func summarizeV2Segments(segs []v2Segment) string {
	parts := make([]string, 0, len(segs))
	for _, seg := range segs {
		parts = append(parts, summarizeV2Segment(seg))
	}
	return strings.Join(parts, "+")
}

func summarizeV2Segment(seg v2Segment) string {
	switch s := seg.(type) {
	case v2Literal:
		return "Literal"
	case v2PlaceholderRef:
		return fmt.Sprintf("Placeholder{%s}", s.Name)
	case v2Color:
		return "Color " + strings.Join(s.Spec.Tokens, " ")
	default:
		return "Unknown"
	}
}

func combineV2SummaryParts(parts []string) string {
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