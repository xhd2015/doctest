// Virtual package body: appears as example.com/app/__doctest_internal_expose/greet
// via go -overlay. Lives under product module path tree so it may import internal.
package expose

import src "example.com/app/internal/greet"

// Hello re-exports the product internal API for external consumers (e.g. mapping-gen).
func Hello() string { return src.Hello() }
