# Scenario

**Feature**: leaf Setup reads `DOCTEST_CACHE_ENV_PROBE` via `os.Getenv` — env not in key

## Preconditions
- Parent `env-getenv` harness runs warmup + two captured runs (probe A then B).
- Leaf-cache key is spine-only; getenv values are not hashed.

## Steps
1. Configure leaf Setup to call `os.Getenv("DOCTEST_CACHE_ENV_PROBE")`.
2. Run the env-cache sequence; both measured runs should stay Cached.

```go
import (
    "testing"
)

func Setup(t *testing.T, req *Request) error {
    cfg := envCacheCfg{
        LeafSetupGo: `import (
    "os"
    "testing"
)

func Setup(t *testing.T, req *Request) error {
    _ = os.Getenv("DOCTEST_CACHE_ENV_PROBE")
    return nil
}`,
    }
    doEnvCacheRun(t, req, cfg)
    return nil
}
```
