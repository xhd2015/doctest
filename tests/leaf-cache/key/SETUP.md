# Scenario

**Feature**: ComputeLeafKey over the local content DAG

```
# leaf key from spine + local packages + module identity + go version
moduleRoot + treeRoot + leafDir + goVersion
  -> ComputeLeafKey
  -> hex digest
```

## Preconditions

- Leaves under this node exercise `ComputeLeafKey` only (no Store).
- Default flavor is `base` unless a child selects `replace` or `remote`.

## Steps

1. Build a mini module + doctest tree via `ensureWorkspace`.
2. Run compute ops (`compute_twice`, `compute_mutate`, or `compute_go_versions`).
3. Assert key equality, inequality, and hex shape.

## Context

- Significance: key computation is the primary product of P1; store is a thin disk map.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	t.Helper()
	if req.Flavor == "" {
		req.Flavor = "base"
	}
	if req.GoVersion == "" {
		req.GoVersion = "go1.25.0"
	}
	return ensureWorkspace(t, req)
}
```
