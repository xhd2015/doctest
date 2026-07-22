# Scenario

**Feature**: workspace skip env is tree-qualified — warm tree-a does not skip tree-b at same relpath

```
run1: test tree-a -> exit 0, store tree-a
run2: test <mod>/... -> exit != 0 (tree-b fail body), Cached == 1
```

## Preconditions

- Parent prepared same-relpath pass/fail workspace and Args/Args2.

## Steps

1. Keep configuration.
2. Assert run1 exit 0; run2 exit != 0 and Cached == 1.

## Context

- If Cached == 2 and exit 0, bare relpath skip falsely covered tree-b.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	t.Helper()
	req.Op = "runtime_multi"
	return nil
}
```
