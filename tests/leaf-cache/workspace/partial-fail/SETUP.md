# Scenario

**Feature**: workspace partial fail PutPass only for non-failed leaves; warm re-run skips them

```
run1: tree-a pass + tree-b fail -> exit != 0; PutPass only a_pass
run2: a_pass Cached/skip; b_fail re-executes -> exit != 0, Cached >= 1
```

## Preconditions

- Mixed multi-tree fixture (`prepareWorkspacePartialFail`).
- Leaf-cache enabled (no disable flags).

## Steps

1. Build partial-fail workspace module.
2. Args = Args2 = `test <mod>/...`.
3. Assert both exits non-zero; run2 Cached >= 1.

## Context

- Complements library `record-partial-fail` with product workspace path.
- If run2 shows 0 Cached, RecordPasses was not wired for workspace or fail
  incorrectly prevented storing the pass leaf.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	t.Helper()
	req.Op = "runtime_multi"
	prepareWorkspacePartialFail(t, req)
	pat := mustWorkspacePattern(req.WorkDir)
	req.Args = []string{"test", pat}
	req.Args2 = []string{"test", pat}
	return nil
}
```
