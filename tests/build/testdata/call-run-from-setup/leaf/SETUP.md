# Scenario

**Feature**: tests for leaf

```
# parse test tree, generate Go code, compile binary
doctest build <test-dir> -> .md files -> Go code -> go build -> binary

# gen-dir controls output layout
gen-dir -> per-leaf packages -> file system
```

## Steps

- Call Run(t, d, req) from inside Setup.
- The generator emits Run as lowercase `run`, then aliases `Run := run`,
  so uppercase `Run` in the Setup body is the harness Run (with inject `d`).

```go
import "github.com/xhd2015/doctest/session"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
    req.Name = "leaf"
    resp, runErr := Run(t, d, req)
    _ = resp
    _ = runErr
    return nil
}
```
