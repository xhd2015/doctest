package spec

import (
	"strings"
	"testing"
)

func TestContentKnownSkills(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{name: "doc-spec", want: "doc-style-test-specification"},
		{name: "code-spec", want: "doc-style-test-code-specification"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			content, err := Content(tt.name)
			if err != nil {
				t.Fatalf("Content(%q): %v", tt.name, err)
			}
			if !strings.Contains(content, tt.want) {
				t.Fatalf("content missing %q", tt.want)
			}
		})
	}
}

func TestContentUnknownSkill(t *testing.T) {
	_, err := Content("unknown")
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "unknown skill: unknown") {
		t.Fatalf("unexpected error: %v", err)
	}
}
