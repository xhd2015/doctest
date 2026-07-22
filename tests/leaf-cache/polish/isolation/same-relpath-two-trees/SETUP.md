# Scenario

**Feature**: same relative leaf path under different TreeRoots → distinct keys

```
treeA/leaf + treeB/leaf (identical content)
  -> ComputeLeafKey(A) != ComputeLeafKey(B)
```

## Preconditions

- Content-identical twins; relative path `leaf` in both.
- Op=`compute_two_inputs`.

## Steps

1. `prepareTwinTrees`.
2. Run two independent KeyInputs.

## Context

- Library tree-identity (also covered under `key/tree-identity/`); product
  cross-tree warm skip remains L3 under `workspace/isolation/**` and was
  formerly e2e here — demoted to L2 key proof for discovery speed.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	t.Helper()
	req.Op = "compute_two_inputs"
	prepareTwinTrees(t, req)
	return nil
}
```
