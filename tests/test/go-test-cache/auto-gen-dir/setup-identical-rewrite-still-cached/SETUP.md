# Scenario

**Feature**: a first run has completed successfully

```
# build and run test binary, report results
doctest test <dir> -> build -> run binary -> pass/fail per leaf -> exit code

# path patterns
.../ -> walk tree | subdir -> run subtree | multi-dir -> aggregate results

# output
progress dots -> . F | verbose -> go test -v | count -> N tests
```

## Preconditions
- A first run has completed successfully.
- The auto-gen-dir parent provides multi-run Run.

## Steps
1. Rewrite leaf SETUP.md with byte-identical content (no semantic change).

```go
import (
    "testing"
)

func Setup(t *testing.T, req *Request) error {
    cfg := multiRunCfg{
        TestDir:          createTempTestProject(t, "mytest"),
        ModifyFile:       "simple/SETUP.md",
        RewriteIdentical: true,
    }
    doMultiRun(t, req, cfg)
    return nil
}
```
