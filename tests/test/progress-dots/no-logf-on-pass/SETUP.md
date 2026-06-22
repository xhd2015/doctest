# Scenario

**Feature**: non-verbose `doctest test` must not forward `t.Logf` output from generated tests

```
# non-verbose run
doctest test <dir> -> build.Test (no -v) -> go test -json -> dots + summary only

# regression: -json parsing must not print t.Logf lines on pass
generated Setup calls t.Logf -> must not appear on stdout or stderr
```

## Preconditions
- A passing leaf whose `Setup` calls `t.Logf` with a unique marker string.

## Steps
1. Create the logf pass tree via `createLogfPassTree`.
2. Run `doctest test <dir>` without `-v`.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
    testDir := createLogfPassTree(t)
    req.Args = []string{"test", testDir}
    return nil
}
```