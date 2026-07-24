# Scenario

**Feature**: any `-count` (including `-count=1`) bypasses leaf-cache → 0 Cached

```
# warm: default flags store leaf-cache pass
doctest test <fixture> -> Cached > 0 on subsequent default run

# -count=1: leaf-cache skip disabled (matches go: -count is non-cacheable)
doctest test <fixture> -count=1 -> 0 Cached
```

## Preconditions
- A first measured run completes successfully with leaf-cache warm hit.
- The auto-gen-dir parent provides multi-run harness.

## Steps
1. Set cfg.SecondFlags to `-count=1` so the second measured run bypasses leaf-cache.

```go
import (
    "testing"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
    cfg := multiRunCfg{}
    cfg.TestDir = createTempTestProject(t, "mytest")
    cfg.SecondFlags = []string{"-count=1"}
    doMultiRun(t, req, cfg)
    return nil
}
```
