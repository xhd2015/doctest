package validate

import (
	"fmt"
	"strings"
)

func hasDSNSection(content string) bool {
	return strings.Contains(content, "DSN (Domain Specific Notion)")
}

func hasVersionSection(content string) bool {
	return strings.Contains(content, "## Version")
}

func hasScenarioSection(content string) bool {
	trimmed := strings.TrimLeft(content, " \t")
	if i := strings.Index(trimmed, "\n"); i >= 0 {
		trimmed = trimmed[:i]
	}
	return trimmed == "# Scenario"
}

func checkDOCTESTSections(path, content string) error {
	if !hasDSNSection(content) {
		return fmt.Errorf("%s: DOCTEST.md must include a DSN (Domain Specific Notion) section", path)
	}
	if !hasVersionSection(content) {
		return fmt.Errorf("%s: DOCTEST.md must include a ## Version section", path)
	}
	return nil
}

func checkSETUPSections(path, content string) error {
	if hasScenarioSection(content) {
		return nil
	}
	return fmt.Errorf("%s: SETUP.md must start with a # Scenario section", path)
}