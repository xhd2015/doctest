package assert

import (
	"fmt"
	"sort"
	"strings"
)

// v3PatternSummary returns a stable summary string for v3 pattern introspection.
func v3PatternSummary(p v3Pattern) string {
	var parts []string
	if len(p.placeholders) > 0 {
		parts = append(parts, v3SummarizePlaceholders(p.placeholders))
	}
	for _, it := range p.items {
		parts = append(parts, v3SummarizeItem(it))
	}
	return v3CombineSummaryParts(parts)
}

func v3SummarizePlaceholders(placeholders map[string]v3Placeholder) string {
	names := make([]string, 0, len(placeholders))
	for _, ph := range placeholders {
		names = append(names, ph.Name)
	}
	sort.Strings(names)
	var defs []string
	for _, name := range names {
		for _, ph := range placeholders {
			if ph.Name == name {
				defs = append(defs, v3SummarizePlaceholder(ph))
				break
			}
		}
	}
	return "Placeholders:" + strings.Join(defs, ",")
}

func v3SummarizePlaceholder(ph v3Placeholder) string {
	s := fmt.Sprintf("%s{%s", ph.Name, ph.Type)
	keys := make([]string, 0, len(ph.Metadata))
	for k := range ph.Metadata {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		s += "," + k + "=" + ph.Metadata[k]
	}
	if ph.Type == "regex" && ph.Regex != "" {
		s += ",regex=" + ph.Regex
	}
	s += "}"
	return s
}

func v3SummarizeItem(it v3Item) string {
	switch x := it.(type) {
	case v3RegexLine:
		return "RegexLine"
	case v3OmitLine:
		return fmt.Sprintf("OmitLine{%d}", x.Count)
	default:
		return "Unknown"
	}
}

func v3CombineSummaryParts(parts []string) string {
	if len(parts) == 0 {
		return ""
	}
	return strings.Join(parts, "+")
}
