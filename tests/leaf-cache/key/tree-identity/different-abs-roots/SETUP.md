# Scenario

**Feature**: different absolute TreeRoot values yield different keys for identical content

```
treeA/leaf (content C) -> keyA
treeB/leaf (content C) -> keyB
keyA != keyB
```

## Preconditions

- Twin trees from parent Setup.
- Op already `compute_two_inputs`.

## Steps

1. Compute both keys.
2. Assert both hex and unequal.

## Context

- Classic TDD for product fix: mix abs TreeRoot (or stable tree id) into ComputeLeafKey.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	t.Helper()
	req.Op = "compute_two_inputs"
	return nil
}
```
