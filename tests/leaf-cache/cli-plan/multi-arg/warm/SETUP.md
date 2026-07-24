# Scenario

**Feature**: multi-arg warm path stores passes and skips on second run

```
run1: doctest test tree-a tree-b (default) -> all pass, PutPass each leaf
run2: doctest test tree-a tree-b (default) -> Cached across both trees
```

## Preconditions

- Parent prepared all-pass multi-tree fixture and multi-arg Args/Args2.
- No `-count` / `-a` on either run.
- Fresh GOCACHE per invocation (parent `runtime_multi`).

## Steps

1. Keep double-run multi-arg configuration.
2. Assert both exit 0; run2 total Cached >= 2 (both tree leaves warm).

## Context

- Primary P3 exit criterion: multi-arg spanning two trees respects leaf-cache
  when skip enabled (same product policy as workspace warm).
- Use **sum** of all `N Cached` summary lines so N× per-tree summaries and a
  future single aggregated summary both pass when both leaves warm-skip.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	t.Helper()
	req.Op = "runtime_multi"
	return nil
}
```
