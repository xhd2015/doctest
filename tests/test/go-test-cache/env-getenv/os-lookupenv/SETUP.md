# Scenario

**Feature**: generated leaf reads `DOCTEST_SESSION_ID` via `os.LookupEnv`

## Preconditions
- Parent `env-getenv` harness runs warmup + two captured runs.

## Steps
1. Configure leaf Setup to call `os.LookupEnv("DOCTEST_SESSION_ID")`.
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
    _, _ = os.LookupEnv("DOCTEST_SESSION_ID")
    return nil
}`
    doEnvCacheRun(t, req)
    return nil
}
```