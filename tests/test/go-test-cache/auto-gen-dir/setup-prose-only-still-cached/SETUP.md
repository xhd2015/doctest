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
1. Rewrite leaf SETUP.md with different markdown prose but **identical** Go block.

```go
import (
    "testing"
)

func Setup(t *testing.T, req *Request) error {
    // Same Go as leafSetupContent(); only Steps prose differs. Build fences via bt.
    proseOnly := "## Steps\n1. prose-only change that must not alter extracted Go\n\n" +
        bt + "go\nimport \"testing\"\nfunc Setup(t *testing.T, req *Request) error { _ = req; return nil }\n" + bt + "\n"
    cfg = multiRunCfg{
        TestDir:       createTempTestProject(t, "mytest"),
        ModifyFile:    "simple/SETUP.md",
        ModifyContent: proseOnly,
    }
    doMultiRun(t, req)
    return nil
}
```
