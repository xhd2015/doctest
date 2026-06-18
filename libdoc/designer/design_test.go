package designer

import (
	"strings"
	"testing"
)

func TestPromptContentResolvesDesignSpec(t *testing.T) {
	resolved := PromptContent()
	for _, placeholder := range []string{"__DOCTEST_SPEC__", "__DOCTEST_DESIGN_SPEC__"} {
		if strings.Contains(resolved, placeholder) {
			t.Fatalf("PromptContent() still contains unresolved placeholder %q", placeholder)
		}
	}
	for _, want := range []string{
		"DSN (Domain Specific Notion)",
		"# Scenario",
	} {
		if !strings.Contains(resolved, want) {
			t.Fatalf("PromptContent() missing %q", want)
		}
	}
}

func TestRawPromptContentHasUnresolvedPlaceholders(t *testing.T) {
	if !strings.Contains(promptContent, "__DOCTEST_DESIGN_SPEC__") {
		t.Fatal("embedded PROMPT.md should contain __DOCTEST_DESIGN_SPEC__ placeholder")
	}
	if !strings.Contains(promptContent, "__DOCTEST_SPEC__") {
		t.Fatal("embedded PROMPT.md should contain __DOCTEST_SPEC__ placeholder")
	}
}