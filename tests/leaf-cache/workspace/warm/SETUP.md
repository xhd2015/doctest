# Scenario

**Feature**: multi-tree workspace warm path stores passes and skips on second run

```
run1: doctest test <mod>/... (default) -> all pass, PutPass each leaf
run2: doctest test <mod>/... (default) -> Cached > 0 (leaf-cache skips)
```

## Preconditions

- Two-tree all-pass fixture (`prepareWorkspaceAllPass`).
- No `-count` / `-a` on either run.
- Fresh GOCACHE per invocation (parent `runtime_multi`).

## Steps

1. Build all-pass workspace module.
2. Args = Args2 = `test <mod>/...`.
3. Assert both exit 0; run2 Cached > 0 (prefer Cached >= 2 when both trees warm).

## Context

- Primary exit criterion for P2 workspace leaf-cache wiring.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	t.Helper()
	req.Op = "runtime_multi"
	prepareWorkspaceAllPass(t, req)
	pat := mustWorkspacePattern(req.WorkDir)
	req.Args = []string{"test", pat}
	req.Args2 = []string{"test", pat}
	return nil
}
```
