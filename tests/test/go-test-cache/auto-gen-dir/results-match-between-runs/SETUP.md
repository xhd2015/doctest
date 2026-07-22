# Scenario

**Feature**: two consecutive runs on the same test tree produce identical test results

```
# build and run test binary, report results
doctest test <dir> -> build -> run binary -> pass/fail per leaf -> exit code

# path patterns
.../ -> walk tree | subdir -> run subtree | multi-dir -> aggregate results

# output
progress dots -> . F | verbose -> go test -v | count -> N tests
```

## Preconditions
- Two consecutive runs on the same test tree produce identical test results.
- The auto-gen-dir parent provides multi-run Run.

## Steps
1. Use parent defaults: two runs, no modifications.

```go
import (
    "testing"
)

func Setup(t *testing.T, req *Request) error {
    cfg := multiRunCfg{}
    cfg.TestDir = createTempTestProject(t, "mytest")
    doMultiRun(t, req, cfg)
    return nil
}
```
