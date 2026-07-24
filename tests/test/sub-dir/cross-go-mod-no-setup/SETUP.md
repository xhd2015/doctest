# Scenario

**Feature**: a module with go.mod but no SETUP.md anywhere

```
# build and run test binary, report results
doctest test <dir> -> build -> run binary -> pass/fail per leaf -> exit code

# path patterns
.../ -> walk tree | subdir -> run subtree | multi-dir -> aggregate results

# output
progress dots -> . F | verbose -> go test -v | count -> N tests
```

## Preconditions
- A module with go.mod but no SETUP.md anywhere.
- No DOCTEST.md in the tree.

## Steps
1. Create mod-c with go.mod and an empty sub-dir.
2. Run `doctest test <mod-c>/empty`.

```go
import (
    "os"
    "path/filepath"
    "testing"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
    tmp := t.TempDir()

    modC := filepath.Join(tmp, "mod-c")
    os.MkdirAll(modC, 0755)
    os.WriteFile(filepath.Join(modC, "go.mod"), []byte("module mod-c\n\ngo 1.21\n"), 0644)

    emptyDir := filepath.Join(modC, "empty")
    os.MkdirAll(emptyDir, 0755)

    req.Args = []string{"test", emptyDir}
    return nil
}
```
