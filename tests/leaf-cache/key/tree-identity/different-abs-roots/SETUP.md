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

- Store keys mix abs TreeRoot into ComputeLeafKey (GREEN). Multi-prep **identity**
  tokens (skip/fail maps) are covered under `runsuite/identity/` separately.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	t.Helper()
	req.Op = "compute_two_inputs"
	return nil
}
```
