# Scenario

**Feature**: `-count=1` disables programmatic leaf-cache skip on multi-tree workspace

```
run1: test <mod>/...            -> store passes
run2: test <mod>/...            -> Cached > 0 (warm works)
run3: test <mod>/... -count=1   -> 0 Cached (bodies re-run)
```

## Preconditions

- All-pass two-tree workspace fixture.
- Three-run sequence: prove warm first so run3's 0 Cached is not vacuous.

## Steps

1. Prepare all-pass workspace; Args/Args2 default `/...`.
2. Child sets Args3 with `-count=1`.
3. Assert run2 Cached > 0; run3 Cached == 0; all exits 0.

## Context

- Aligns workspace path with single-tree `runtime/disable/count-bypasses`.
- Only `-count` is required for P2; `-a` / `--no-leaf-cache` already sealed
  on single-tree (same SkipEnabled helper).

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	t.Helper()
	req.Op = "runtime_multi"
	prepareWorkspaceAllPass(t, req)
	pat := mustWorkspacePattern(req.WorkDir)
	req.Args = []string{"test", pat}
	req.Args2 = []string{"test", pat}
	// Args3 set by child
	return nil
}
```
