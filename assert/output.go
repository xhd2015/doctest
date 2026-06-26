package assert

import "testing"

// Output parses template, matches actual output, and calls t.Fatal on failure.
func Output(t *testing.T, actual, template string, opts ...Option) {
	t.Helper()
	p, err := Parse(template)
	if err != nil {
		t.Fatalf("parse output template: %v", err)
	}
	if err := Match(p, actual, opts...); err != nil {
		t.Fatal(err)
	}
}