# Scenario

**Feature**: quiet + color progress dots distinguish warm leaf-cache skips (grey) from executed pass (plain) and fail (red)

```
# warm second run (quiet, --color)
doctest test --color fixture
  -> identity in this-run skip set -> progress <ansi gray .>
  -> executed pass -> plain .
  -> fail -> <ansi red .>
# -count bypass
doctest test --color -count=1 fixture -> 0 Cached; no grey leaf-cache skip dots
```

## Preconditions

- Nested CLI under `runtime/**`; `--color` forces ANSI when stdout is a pipe.
- Multi-leaf fixtures so progress region has multiple dots.
- Fresh GOCACHE per run; isolated leaf-cache store.

## Steps

1. Parent runtime sets Bin, timeout, isolateRuntimeEnv.
2. Child leaves choose warm / fail / -count scenarios with `--color`.
3. Assert grey / red / absent grey dots in the progress region (before summary).

## Context

- Summary **Cached** segment is already grey when color is on (out of scope here).
- This branch locks **progress** dots only. Shared JSON consumer path covers
  single-tree and workspace/`./...`.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	t.Helper()
	req.Op = "runtime_multi"
	return nil
}
```
