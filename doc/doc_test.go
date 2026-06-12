package doc

import (
	"strings"
	"testing"
)

func TestContent_DocStyleTestSpecification(t *testing.T) {
	content, err := Content("DOC_STYLE_TEST_SPECIFICATION.md")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if content == "" {
		t.Fatal("expected non-empty content")
	}
	if !strings.HasPrefix(content, "---") {
		t.Fatalf("expected content to start with YAML frontmatter '---', got: %q", content[:min(len(content), 50)])
	}
}

func TestContent_DocStyleTestCodeSpecification(t *testing.T) {
	content, err := Content("DOC_STYLE_TEST_CODE_SPECIFICATION.md")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if content == "" {
		t.Fatal("expected non-empty content")
	}
	if !strings.HasPrefix(content, "---") {
		t.Fatalf("expected content to start with YAML frontmatter '---', got: %q", content[:min(len(content), 50)])
	}
}

func TestContent_UnknownFile(t *testing.T) {
	_, err := Content("nonexistent.md")
	if err == nil {
		t.Fatal("expected error for unknown file, got nil")
	}
}
