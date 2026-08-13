# Scenario

**Bug**: generated suite does not compile when a SETUP helper uses `goto` over a later `:=`

```
# wrk unwind-pipeline installFakeOpencode: goto mock jumps over cmd :=
# doctest test testdata/goto-jumps-over-decl -> go test compile
# today: "goto mock jumps over declaration of cmd"
# desired: generated suite compiles; inner leaf runs (exit 0)
```

## Preconditions

- Testdata tree `build/testdata/goto-jumps-over-decl` has a helper with
  `goto mock` before `cmd := exec.Command(...)`.
- Same pattern as wrk `cmd/wrk/tests/unwind-pipeline` SETUP.

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
