# Scenario

**Feature**: no target directory argument is supplied

```
# parse test tree, generate Go code, compile binary
doctest build <test-dir> -> .md files -> Go code -> go build -> binary

# gen-dir controls output layout
gen-dir -> per-leaf packages -> file system
```

## Preconditions
- No target directory argument is supplied.

## Steps
1. Run `doctest build`.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
    req.Args = []string{"build"}
    return nil
}
```

