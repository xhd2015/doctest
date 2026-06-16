# Scenario

**Feature**: a doc-style test tree where a child SETUP.md's `func Setup` calls `Run(t, req)`

```
# parse test tree, generate Go code, compile binary
doctest build <test-dir> -> .md files -> Go code -> go build -> binary

# gen-dir controls output layout
gen-dir -> per-leaf packages -> file system
```

## Preconditions
- A doc-style test tree where a child SETUP.md's `func Setup` calls `Run(t, req)`.
- The generator lowers `func Run` to lowercase closure `run`, but the Setup body
  still references uppercase `Run` → "undefined: Run".

## Steps
1. Run `doctest test` on the testdata fixture, which compiles with `go test -c`.

```go
import "path/filepath"

func Setup(t *testing.T, req *Request) error {
    req.Args = []string{"test", filepath.Join(DOCTEST_ROOT, "build", "testdata", "call-run-from-setup")}
    return nil
}
```
