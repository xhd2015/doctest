package validate

import "testing"

func TestHasVersionSection(t *testing.T) {
	tests := []struct {
		name    string
		content string
		want    bool
	}{
		{name: "present", content: "## Version\n0.0.2\n", want: true},
		{name: "missing", content: "# Tests\n", want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := hasVersionSection(tt.content); got != tt.want {
				t.Fatalf("hasVersionSection() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHasDSNSection(t *testing.T) {
	tests := []struct {
		name    string
		content string
		want    bool
	}{
		{name: "heading with full title", content: "## DSN (Domain Specific Notion)\n\nparticipants", want: true},
		{name: "missing", content: "# Tests\n\noverview only", want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := hasDSNSection(tt.content); got != tt.want {
				t.Fatalf("hasDSNSection() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHasScenarioSection(t *testing.T) {
	tests := []struct {
		name    string
		content string
		want    bool
	}{
		{name: "first section", content: "# Scenario\n\n**Feature**: x\n", want: true},
		{name: "leading whitespace", content: "  # Scenario\n", want: true},
		{name: "not first", content: "## Preconditions\n\n# Scenario\n", want: false},
		{name: "missing", content: "## Preconditions\n", want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := hasScenarioSection(tt.content); got != tt.want {
				t.Fatalf("hasScenarioSection() = %v, want %v", got, tt.want)
			}
		})
	}
}