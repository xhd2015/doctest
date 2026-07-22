# Scenario

**Feature**: `-a` (force rebuild) bypasses leaf-cache → 0 Cached

```
# warm: default flags store leaf-cache pass
doctest test <fixture> -> Cached > 0 on subsequent default run

# -a: leaf-cache skip disabled (and go test -a rebuild)
doctest test <fixture> -a -> 0 Cached
```

## Preconditions
- A first measured run completes successfully with leaf-cache warm hit.
- The auto-gen-dir parent provides multi-run harness.

## Steps
1. Set cfg.SecondFlags to `-a` so the second measured run bypasses leaf-cache.

```go
import (
    "testing"
)

func Setup(t *testing.T, req *Request) error {
    cfg := multiRunCfg{}
    cfg.TestDir = createTempTestProject(t, "mytest")
    cfg.SecondFlags = []string{"-a"}
    doMultiRun(t, req, cfg)
    return nil
}
```
