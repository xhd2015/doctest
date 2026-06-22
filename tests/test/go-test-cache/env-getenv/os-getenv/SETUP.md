# Scenario

**Feature**: generated leaf reads `DOCTEST_SESSION_ID` via `os.Getenv`

## Preconditions
- Parent `env-getenv` harness runs warmup + two captured runs.

## Steps
1. Configure leaf Setup to call `os.Getenv("DOCTEST_SESSION_ID")`.
2. Run the env-cache proof sequence.

```go
import (
    "testing"
)

func Setup(t *testing.T, req *Request) error {
    envCfg.LeafSetupGo = `import (
    "os"
    "testing"
)

func Setup(t *testing.T, req *Request) error {
    _ = os.Getenv("DOCTEST_SESSION_ID")
    return nil
}`
    doEnvCacheRun(t, req)
    return nil
}
```