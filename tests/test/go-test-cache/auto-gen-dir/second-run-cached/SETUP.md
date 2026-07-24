# Scenario

**Feature**: a temporary test project exists with a valid doctest tree

```
# build and run test binary, report results
doctest test <dir> -> build -> run binary -> pass/fail per leaf -> exit code

# path patterns
.../ -> walk tree | subdir -> run subtree | multi-dir -> aggregate results

# output
progress dots -> . F | verbose -> go test -v | count -> N tests
```

## Preconditions
- A temporary test project exists with a valid doctest tree.
- The auto-gen-dir parent provides multi-run Run.

## Steps
1. Use parent defaults: two runs, no modifications, no count override.

```go
import (
    "testing"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
    cfg := multiRunCfg{}
    cfg.TestDir = createTempTestProject(t, "mytest")
    doMultiRun(t, req, cfg)
    return nil
}
```
