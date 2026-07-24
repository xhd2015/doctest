# Scenario

**Feature**: P3 polish — selective key invalidation (L2) + help docs (L3)

```
# L2 library (no nested doctest):
# selective: edit one leaf ASSERT → sibling key stable
# local dep: edit imported package → leaf key changes
# isolation: twin trees same relpath → distinct keys

# L3 e2e (label: heavy):
# docs: test --help lists leaf-cache flags (product binary)
```

## Preconditions

- **Most leaves are L2** — `ComputeLeafKey` only; no `testbin` at this node.
- **Docs** child builds the product binary for `runtime_once` help.
- Selective / isolation / local-dep are unlabeled (discovery).

## Steps

1. Default GoVersion for library leaves.
2. Child selects fixture + Op (library or `runtime_once` for docs).

## Context

- Key-level selective invalidation and tree isolation live here as L2 mass.
- Full product warm **Cached** wiring remains under L3 `runtime/**`.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	t.Helper()
	if req.GoVersion == "" {
		req.GoVersion = "go1.25.0"
	}
	return nil
}
```
