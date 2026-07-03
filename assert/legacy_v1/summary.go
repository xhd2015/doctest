package legacy_v1

import (
	"fmt"
	"strings"
)

// PatternSummary returns a stable summary string for test introspection.
func PatternSummary(p Pattern) string {
	parts := make([]string, 0, len(p.Items))
	for _, item := range p.Items {
		parts = append(parts, summarizeItem(item))
	}
	return combineSummaryParts(parts)
}

func combineSummaryParts(parts []string) string {
	if len(parts) == 0 {
		return ""
	}
	if len(parts) >= 2 && parts[0] == "LiteralLine" {
		allLiteral := true
		for _, p := range parts {
			if p != "LiteralLine" {
				allLiteral = false
				break
			}
		}
		if allLiteral {
			return fmt.Sprintf("LiteralLine×%d", len(parts))
		}
	}
	return strings.Join(parts, "+")
}

func summarizeItem(item Item) string {
	switch it := item.(type) {
	case LiteralLine:
		return "LiteralLine"
	case PatternLine:
		return "PatternLine segments:" + summarizeSegments(it.Segments)
	case RegexLine:
		return "RegexLine"
	case BlockOptional:
		if len(it.Items) == 0 {
			return "BlockOptional{}"
		}
		return fmt.Sprintf("BlockOptional{%d}", len(it.Items))
	case AnyOfBlock:
		return fmt.Sprintf("AnyOfBlock branches:%d", len(it.Branches))
	case ContainsBlock:
		return summarizeContainsBlock(it)
	default:
		return "Unknown"
	}
}

func summarizeContainsBlock(cb ContainsBlock) string {
	s := fmt.Sprintf("ContainsBlock fragments:%d", len(cb.Fragments))
	for _, f := range cb.Fragments {
		switch f.Mode {
		case ContainsStartWith:
			s += " StartWith"
		case ContainsEndWith:
			s += " EndWith"
		}
	}
	return s
}

func summarizeSegments(segs []Segment) string {
	parts := make([]string, 0, len(segs))
	for _, seg := range segs {
		parts = append(parts, summarizeSegment(seg))
	}
	return strings.Join(parts, "+")
}

func literalTagSnippet(text string) string {
	start := strings.Index(text, "<")
	if start < 0 {
		return ""
	}
	end := strings.Index(text[start:], ">")
	if end < 0 {
		return ""
	}
	return text[start : start+end+1]
}

func summarizeSegment(seg Segment) string {
	switch s := seg.(type) {
	case Literal:
		if tag := literalTagSnippet(s.Text); tag != "" {
			return "Literal " + tag
		}
		return "Literal"
	case Hint:
		return "Hint:" + s.Label
	case InlineOptional:
		return "InlineOptional"
	case InlineAnyOf:
		return fmt.Sprintf("InlineAnyOf branches:%d", len(s.Branches))
	case AnsiColor:
		return "AnsiColor " + strings.Join(s.Spec.Tokens, " ")
	case InlineRegex:
		return "InlineRegex"
	default:
		return "Unknown"
	}
}