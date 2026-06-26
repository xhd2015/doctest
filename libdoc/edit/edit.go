package edit

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/xhd2015/doctest/libdoc/core"
)

// Edit updates ASSERT.md frontmatter for a single concrete leaf test.
func Edit(path string, addLabel string, addExplanation string) error {
	if strings.Contains(path, "...") {
		return fmt.Errorf("edit does not accept ... patterns")
	}
	abs, err := filepath.Abs(path)
	if err != nil {
		return err
	}
	leafDir, assertPath, err := resolveLeafAssert(abs)
	if err != nil {
		return err
	}
	if addLabel == "" && addExplanation == "" {
		return fmt.Errorf("edit requires --add-label and/or --add-explanation")
	}

	content, err := os.ReadFile(assertPath)
	if err != nil {
		return err
	}
	fm, body, err := core.ParseAssertFrontmatter(assertPath, string(content))
	if err != nil {
		return err
	}

	if addLabel != "" {
		label := strings.TrimSpace(addLabel)
		if label == "" {
			return fmt.Errorf("--add-label requires a non-empty value")
		}
		if containsLabel(fm.Labels, label) {
			fmt.Fprintf(os.Stderr, "warning: label %q already set on %s\n", label, leafDir)
		} else {
			fm.Labels = append(fm.Labels, label)
		}
	}
	if addExplanation != "" {
		text := strings.TrimSpace(addExplanation)
		if fm.Explanation == "" {
			fm.Explanation = text
		} else {
			fm.Explanation = fm.Explanation + "; " + text
		}
	}

	updated := renderAssertWithFrontmatter(fm, body)
	if err := os.WriteFile(assertPath, []byte(updated), 0644); err != nil {
		return err
	}
	return nil
}

func resolveLeafAssert(abs string) (leafDir string, assertPath string, err error) {
	if filepath.Base(abs) == "ASSERT.md" {
		return filepath.Dir(abs), abs, nil
	}
	info, statErr := os.Stat(abs)
	if statErr != nil {
		return "", "", statErr
	}
	if !info.IsDir() {
		return "", "", fmt.Errorf("%s: not a doctest leaf (expected directory with ASSERT.md or ASSERT.md path)", abs)
	}
	assertPath = filepath.Join(abs, "ASSERT.md")
	if _, err := os.Stat(assertPath); err != nil {
		return "", "", fmt.Errorf("%s: not a doctest leaf (ASSERT.md missing)", abs)
	}
	return abs, assertPath, nil
}

func containsLabel(labels []string, label string) bool {
	for _, l := range labels {
		if l == label {
			return true
		}
	}
	return false
}

func renderAssertWithFrontmatter(fm core.AssertFrontmatter, body string) string {
	if len(fm.Labels) == 0 && fm.Explanation == "" {
		return strings.TrimLeft(body, "\r\n")
	}
	var b strings.Builder
	b.WriteString("---\n")
	if len(fm.Labels) > 0 {
		b.WriteString("label: ")
		b.WriteString(strings.Join(fm.Labels, ", "))
		b.WriteString("\n")
	}
	if fm.Explanation != "" {
		b.WriteString("explanation: ")
		b.WriteString(fm.Explanation)
		b.WriteString("\n")
	}
	b.WriteString("---\n\n")
	body = strings.TrimLeft(body, "\r\n")
	b.WriteString(body)
	if !strings.HasSuffix(b.String(), "\n") {
		b.WriteString("\n")
	}
	return b.String()
}