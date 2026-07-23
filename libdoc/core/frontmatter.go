package core

import (
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

// AssertFrontmatter is optional YAML metadata at the top of ASSERT.md.
type AssertFrontmatter struct {
	Labels      []string
	Explanation string
}

type assertFrontmatterYAML struct {
	Label       string `yaml:"label"`
	Explanation string `yaml:"explanation"`
}

// ParseAssertFrontmatter parses an optional YAML frontmatter block from ASSERT.md.
// Missing frontmatter is valid. Malformed frontmatter returns an error.
func ParseAssertFrontmatter(path string, content string) (AssertFrontmatter, string, error) {
	yamlText, body, found, err := splitYAMLFrontmatter(content)
	if err != nil {
		return AssertFrontmatter{}, content, fmt.Errorf("%s: %w", path, err)
	}
	if !found {
		return AssertFrontmatter{}, content, nil
	}
	var raw assertFrontmatterYAML
	if err := yaml.Unmarshal([]byte(yamlText), &raw); err != nil {
		return AssertFrontmatter{}, content, fmt.Errorf("%s: invalid frontmatter: %w", path, err)
	}
	return AssertFrontmatter{
		Labels:      parseLabelField(raw.Label),
		Explanation: strings.TrimSpace(raw.Explanation),
	}, body, nil
}

func splitYAMLFrontmatter(content string) (yamlText string, body string, found bool, err error) {
	if !strings.HasPrefix(content, "---") {
		return "", content, false, nil
	}
	rest := content[len("---"):]
	switch {
	case strings.HasPrefix(rest, "\n"):
		rest = rest[1:]
	case strings.HasPrefix(rest, "\r\n"):
		rest = rest[2:]
	default:
		return "", content, true, fmt.Errorf("frontmatter must start with --- followed by a newline")
	}

	closeIdx := indexFrontmatterClose(rest)
	if closeIdx < 0 {
		return "", content, true, fmt.Errorf("frontmatter missing closing ---")
	}
	yamlText = strings.TrimSpace(rest[:closeIdx])
	body = rest[closeIdx:]
	body = strings.TrimLeft(body, "\r\n")
	return yamlText, body, true, nil
}

func indexFrontmatterClose(rest string) int {
	for i := 0; i < len(rest); i++ {
		if rest[i] != '\n' {
			continue
		}
		lineStart := i + 1
		if lineStart+3 > len(rest) {
			continue
		}
		line := rest[lineStart:]
		if strings.HasPrefix(line, "---") {
			after := line[3:]
			if after == "" || after[0] == '\n' || (len(after) >= 2 && after[0] == '\r' && after[1] == '\n') {
				return lineStart + 3
			}
		}
	}
	return -1
}

func parseLabelField(label string) []string {
	label = strings.TrimSpace(label)
	if label == "" {
		return nil
	}
	parts := strings.Split(label, ",")
	var out []string
	seen := make(map[string]bool)
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" || seen[p] {
			continue
		}
		seen[p] = true
		out = append(out, p)
	}
	return out
}

func (tc TreeCase) HasLabel() bool {
	return len(tc.Labels) > 0
}

// PartitionLabeledCases splits runnable leaves for discovery-mode runs.
// When skipLabeled is false (explicit leaf), all cases are returned to run.
func PartitionLabeledCases(cases []TreeCase, skipLabeled bool) (run []TreeCase, skipped []SkippedCase) {
	if !skipLabeled {
		return cases, nil
	}
	for _, tc := range cases {
		if tc.HasLabel() {
			skipped = append(skipped, SkippedCase{
				Name:        tc.Name,
				Path:        tc.Path,
				Labels:      append([]string(nil), tc.Labels...),
				Explanation: tc.Explanation,
			})
			continue
		}
		run = append(run, tc)
	}
	return run, skipped
}

// FilterCasesByLabel applies discovery skip or --label filtering.
// When opts.LabelAll is true, every case runs (labels ignored for selection).
// When len(opts.LabelExprs)==0, behavior matches PartitionLabeledCases with skipLabeled=!opts.ExplicitLeaf.
// When label expressions are set, only leaves whose label set matches any EXPR run
// (boolean eval; unlabeled = empty set, so e.g. !e2e matches). Non-matches are
// skipped with Reason "label filter".
func FilterCasesByLabel(cases []TreeCase, opts Options) (run []TreeCase, skipped []SkippedCase) {
	if opts.LabelAll {
		return cases, nil
	}
	if len(opts.LabelExprs) == 0 {
		return PartitionLabeledCases(cases, !opts.ExplicitLeaf)
	}
	for _, tc := range cases {
		match, err := MatchLabelExprs(opts.LabelExprs, tc.Labels)
		if err != nil {
			// Caller should validate expressions before discovery; treat as non-match if called anyway.
			match = false
		}
		if match {
			run = append(run, tc)
			continue
		}
		skipped = append(skipped, SkippedCase{
			Name:        tc.Name,
			Path:        tc.Path,
			Labels:      append([]string(nil), tc.Labels...),
			Explanation: tc.Explanation,
			Reason:      "label filter",
		})
	}
	return run, skipped
}