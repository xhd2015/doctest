# Scenario

**Feature**: a plain directory is not part of any doc-style test tree (no DOCTEST.md, no SETUP.md anywhere)

```
# build and run test binary, report results
doctest test <dir> -> build -> run binary -> pass/fail per leaf -> exit code

# path patterns
.../ -> walk tree | subdir -> run subtree | multi-dir -> aggregate results

# output
progress dots -> . F | verbose -> go test -v | count -> N tests
```

## Preconditions
- A plain directory is not part of any doc-style test tree (no DOCTEST.md, no SETUP.md anywhere).

## Steps
1. Create an empty plain directory.
2. Run `doctest test <plainDir>`.

```go
import (
    "path/filepath"
    "testing"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
    plainDir := t.TempDir()
    req.Args = []string{"test", plainDir}
    return nil
}
```
