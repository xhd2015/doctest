package leafshim

import "example.com/realmod/http/internal/leaf"

// Bridge re-exports internal Hello so importers outside http/internal can call it.
func Bridge() string { return leaf.Hello() }
