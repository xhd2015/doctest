# Scenario

**Feature**: the test data directory exists at `DOCTEST_ROOT/testdata/basic-request-runner`

```
# build and run test binary, report results
doctest test <dir> -> build -> run binary -> pass/fail per leaf -> exit code

# output
progress dots -> . F | verbose -> go test -v | count -> N tests | timeout -> go test -timeout=...
```

## Preconditions
- The test data directory exists at `DOCTEST_ROOT/testdata/basic-request-runner`.

## Steps
1. Run `doctest test` with `--timeout` and check the printed go test command line.

```go
import (
"github.com/xhd2015/doctest/session"
    "os"
    "path/filepath"
    "testing"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
    exampleDir := filepath.Join(d.DOCTEST_ROOT, "testdata", "basic-request-runner")
    if info, err := os.Stat(exampleDir); err != nil {
        t.Fatalf("testdata dir %s not found: %v", exampleDir, err)
    } else if !info.IsDir() {
        t.Fatalf("testdata dir %s is not a directory", exampleDir)
    }
    return nil
}
```