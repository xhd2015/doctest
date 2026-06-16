# Scenario

**Feature**: the test data directory exists at `DOCTEST_ROOT/testdata/basic-request-runner`

```
# build and run test binary, report results
doctest test <dir> -> build -> run binary -> pass/fail per leaf -> exit code

# path patterns
.../ -> walk tree | subdir -> run subtree | multi-dir -> aggregate results

# output
progress dots -> . F | verbose -> go test -v | count -> N tests
```

## Preconditions
- The test data directory exists at `DOCTEST_ROOT/testdata/basic-request-runner`.

## Steps
1. Run `doctest test` with or without `-v` flag.
2. Check the printed go test command line for `-v` presence.

```go
import (
    "os"
    "path/filepath"
    "testing"
)

func Setup(t *testing.T, req *Request) error {
    exampleDir := filepath.Join(DOCTEST_ROOT, "testdata", "basic-request-runner")
    if info, err := os.Stat(exampleDir); err != nil {
        t.Fatalf("testdata dir %s not found: %v", exampleDir, err)
    } else if !info.IsDir() {
        t.Fatalf("testdata dir %s is not a directory", exampleDir)
    }
    return nil
}
```
