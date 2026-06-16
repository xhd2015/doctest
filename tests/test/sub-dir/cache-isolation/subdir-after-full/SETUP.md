# Scenario

**Feature**: a doc-style test tree with two groups (group-a with leaf-1, leaf-2; group-b with leaf-3)

```
# build and run test binary, report results
doctest test <dir> -> build -> run binary -> pass/fail per leaf -> exit code

# path patterns
.../ -> walk tree | subdir -> run subtree | multi-dir -> aggregate results

# output
progress dots -> . F | verbose -> go test -v | count -> N tests
```

## Preconditions
- A doc-style test tree with two groups (group-a with leaf-1, leaf-2; group-b with leaf-3).
- The parent cache-isolation SETUP.md provides multi-run Run.

## Steps
1. Set scenario to "subdir_after_full".
2. First run executes doctest on the full tree root (caches all tests).
3. Second run executes doctest on the group-a sub-directory only.

```go
import (
    "testing"
)

func Setup(t *testing.T, req *Request) error {
    cisoCfg.Scenario = "subdir_after_full"
    doMultiRun(t, req)
    return nil
}
```
