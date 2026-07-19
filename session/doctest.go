package session

// Doctest is the public inject-contract context for a doctest run.
// Field names match the former free variables of the same names.
// These are struct fields only (not process environment variables).
type Doctest struct {
	DOCTEST_ROOT       string
	DOCTEST_CASE       string
	DOCTEST_SESSION_ID string
}
