# Scenario

**Feature**: `-count=1` disables programmatic leaf-cache skip on multi-arg path

```
run1: test tree-a tree-b            -> store passes
run2: test tree-a tree-b            -> Cached > 0 (warm works)
run3: test tree-a tree-b -count=1   -> 0 Cached (bodies re-run)
```

## Preconditions

- All-pass two-tree multi-arg fixture (parent multi-arg SETUP).
- Three-run sequence: prove warm first so run3's 0 Cached is not vacuous.

## Steps

1. Keep multi-arg Args/Args2 from parent.
2. Child sets Args3 with `-count=1` after the two roots.
3. Assert run2 warm; run3 total Cached == 0; all exits 0.

## Context

- Same SkipEnabled policy as single-tree `runtime/disable/count-bypasses` and
  workspace `workspace/disable/count-bypasses` — multi-arg must not ignore
  `-count` when skip would otherwise apply.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	t.Helper()
	req.Op = "runtime_multi"
	// Args/Args2 already multi-arg from multi-arg/SETUP.md
	return nil
}
```
