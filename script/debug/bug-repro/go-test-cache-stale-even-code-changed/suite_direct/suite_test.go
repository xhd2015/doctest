package suite_direct_test

import (
	"testing"

	"bugrepro.local/gotestcache/mid"
)

// Variant A: test package imports mid directly.
func TestVersion(t *testing.T) {
	if v := mid.Version(); v != 1 {
		t.Fatalf("mid.Version() = %d, want 1", v)
	}
}
