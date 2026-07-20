# Scenario

**Feature**: local replace modules are part of the content DAG

```
# app go.mod: replace example.com/lib => ../lib
# leaf ASSERT imports example.com/lib

lib source change  -> key changes
lib go.mod change  -> key changes
```

## Preconditions

- Flavor = `replace`: sibling `lib/` module + replace directive + import in assert.
- Op is `compute_mutate`.

## Steps

1. Rebuild workspace with replace flavor (override base flavor).
2. Child mutates lib source or lib go.mod.

## Context

- Local replace is how monorepos pull in neighboring packages; both source and
  replace-target go.mod must invalidate.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	t.Helper()
	req.Flavor = "replace"
	req.Op = "compute_mutate"
	// key/ Setup already called ensureWorkspace with base; rebuild with replace.
	req.WorkDir = t.TempDir()
	req.ModuleRoot = ""
	req.TreeRoot = ""
	req.LeafDir = ""
	return ensureWorkspace(t, req)
}
```
