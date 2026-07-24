# Scenario

**Feature**: leaf keys include tree identity so different trees never collide

```
# identical spine content, different absolute TreeRoot
ComputeLeafKey(treeA/leaf) -> keyA
ComputeLeafKey(treeB/leaf) -> keyB
# keyA != keyB
```

## Preconditions

- Two content-identical mini trees under different absolute paths.
- Same GoVersion; no go.mod (module identity empty for both).

## Steps

1. Build twin trees with prepareTwinTrees.
2. Op=`compute_two_inputs`.
3. Assert keys differ.

## Context

- Without tree identity in the hash, relative spine content alone collides across trees.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	t.Helper()
	if req.GoVersion == "" {
		req.GoVersion = "go1.25.0"
	}
	req.Op = "compute_two_inputs"
	prepareTwinTrees(t, req)
	return nil
}
```
