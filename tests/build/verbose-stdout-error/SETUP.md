# Scenario

**Feature**: running `doctest test -v` triggers the verbose code path where

```
# parse test tree, generate Go code, compile binary
doctest build <test-dir> -> .md files -> Go code -> go build -> binary

# gen-dir controls output layout
gen-dir -> per-leaf packages -> file system
```

## Preconditions
- Running `doctest test -v` triggers the verbose code path where
  `goTestBuild.Stdout` and `goTestBuild.Stderr` are set to `w`,
  then `CombinedOutput()` tries to set them again → "Stdout already set".

## Steps
1. Run `doctest test -v` on the testdata fixture.

```go
import (
"path/filepath"
"github.com/xhd2015/doctest/session"
)
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
    req.Args = []string{"test", "-v", filepath.Join(d.DOCTEST_ROOT, "build", "testdata", "verbose-stdout-error")}
    return nil
}
```
