package suite_blank

import (
	"testing"

	// Named import still runs leaf.init (registry-style Fn), like doctest
	// suite depending on leaves only through registration + call.
	"bugrepro.local/gotestcache/leaf"
)

// Variant B: suite reaches mid only through leaf.Fn (function value).
func TestSuite(t *testing.T) {
	if leaf.Fn == nil {
		t.Fatal("leaf.Fn not registered")
	}
	leaf.Fn(t)
}
