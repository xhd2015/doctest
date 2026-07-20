package session

// Metrics is leaf-scoped metrics context for this Run/Setup/Assert invocation.
// Filled by the suite/leaf harness when constructing d — not process environment.
type Metrics struct {
	// ParentLeaf is the outer leaf path that owns nested work (e.g. tree-relative
	// leaf path). Replaces DOCTEST_METRICS_PARENT_LEAF process env.
	ParentLeaf string
	// NestSink is an optional path for nested phase JSONL append.
	// Empty when metrics-off or no nest recording. Prefer explicit field over
	// process env for in-process nest; suite process may still inherit a sink
	// via go test cmd.Env and copy it here at leaf start (read-only).
	NestSink string
}

// Doctest is the public inject-contract context for a doctest run.
// Field names match the former free variables of the same names.
// These are struct fields only (not process environment variables).
type Doctest struct {
	DOCTEST_ROOT       string
	DOCTEST_CASE       string
	DOCTEST_SESSION_ID string
	// Metrics carries nest attribution and sink for this leaf invocation.
	Metrics Metrics
}
