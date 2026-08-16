# Scenario

**Feature**: generated suite compiles a SETUP helper that uses `goto mock`

```
# testdata helper: goto mock; build path := is in an inner block
# doctest test testdata/goto-jumps-over-decl -> go test compile + run
```

## Preconditions

- Testdata tree `build/testdata/goto-jumps-over-decl` has `goto mock` and
  `cmd :=` only inside a nested block (legal Go).

## Steps

1. Run `doctest test` on the testdata fixture (compiles with `go test -c`).

```go
import (
	"path/filepath"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Args = []string{"test", filepath.Join(d.DOCTEST_ROOT, "build", "testdata", "goto-jumps-over-decl")}
	return nil
}
```
