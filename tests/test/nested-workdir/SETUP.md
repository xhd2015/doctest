# Scenario

**Feature**: a valid doc-style test tree exists in the module

```
# build and run test binary, report results
doctest test <dir> -> build -> run binary -> pass/fail per leaf -> exit code

# path patterns
.../ -> walk tree | subdir -> run subtree | multi-dir -> aggregate results

# output
progress dots -> . F | verbose -> go test -v | count -> N tests
```

## Preconditions
- A valid doc-style test tree exists in the module.
- The command is invoked from inside a nested module directory.

## Steps
1. Set the process working directory to the module root (`DOCTEST_ROOT/..`).
2. Run `doctest test <absolute-dir>`.

```go
import (
"github.com/xhd2015/doctest/session"
    "os/exec"
    "path/filepath"
    "testing"

	"github.com/xhd2015/doctest/libdoc/testbin"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
    exampleDir := filepath.Join(d.DOCTEST_ROOT, "testdata/basic-request-runner")
    req.WorkDir = filepath.Join(d.DOCTEST_ROOT, "..")
    req.Args = []string{"test", exampleDir}

    req.Bin = testbin.Ensure(t, filepath.Join(d.DOCTEST_ROOT, ".."))
    return nil
}
```
