package leaf

import (
	"testing"

	"bugrepro.local/gotestcache/mid"
)

// Fn is registered like doctest leaf packages: function value + init.
var Fn func(*testing.T)

func init() {
	Fn = Run
}

// Run exercises mid.Version (same role as leaf calling intermediate Setup).
func Run(t *testing.T) {
	t.Helper()
	if v := mid.Version(); v != 1 {
		t.Fatalf("mid.Version() = %d, want 1", v)
	}
}
