# Scenario

**Feature**: leaf Setup reads `DOCTEST_CACHE_ENV_PROBE` via `syscall.Getenv` — env not in key

## Preconditions
- Parent `env-getenv` harness runs warmup + two captured runs (probe A then B).
- Leaf-cache key is spine-only for all read APIs (including syscall).

## Steps
1. Configure leaf Setup to call `syscall.Getenv("DOCTEST_CACHE_ENV_PROBE")`.
2. Run the env-cache sequence; both measured runs should stay Cached.

```go
import (
    "testing"
)

func Setup(t *testing.T, req *Request) error {
    envCfg.LeafSetupGo = `import (
    "syscall"
    "testing"
)

func Setup(t *testing.T, req *Request) error {
    _, _ = syscall.Getenv("DOCTEST_CACHE_ENV_PROBE")
    return nil
}`
    doEnvCacheRun(t, req)
    return nil
}
```
