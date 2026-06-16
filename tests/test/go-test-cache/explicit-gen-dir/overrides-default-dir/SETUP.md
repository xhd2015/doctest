# Scenario

**Feature**: a valid doctest tree exists

```
# build and run test binary, report results
doctest test <dir> -> build -> run binary -> pass/fail per leaf -> exit code

# path patterns
.../ -> walk tree | subdir -> run subtree | multi-dir -> aggregate results

# output
progress dots -> . F | verbose -> go test -v | count -> N tests
```

## Preconditions
- A valid doctest tree exists.
- The root Setup has built the doctest binary.

## Steps
1. Run `doctest test <dir> --gen-dir <tmp>`.
2. Verify the printed gen dir path equals the specified --gen-dir, not the hash-based path.

```go
import (
    "path/filepath"
    "strings"
    "testing"
)

func Setup(t *testing.T, req *Request) error {
    testDir := createTempTestProject(t, "mytest")
    genDir := filepath.Join(t.TempDir(), "generated")
    req.Args = []string{"test", testDir, "--gen-dir", genDir, "-v"}
    req.Env = append(req.Env, "GEN_DIR="+genDir)
    return nil
}
```
