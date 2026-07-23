# Scenario

**Feature**: internal-compile path resolves assert via temp -modfile instead of nested go.mod

```
# internal/ import detected
doctest -> .doctest_run_* under moduleRoot -> -modfile with parent go.mod + assert replace
```

## Preconditions

- Fixture module imports `example.com/app/internal/greet` in harness `Run()`.
- No nested `go.mod` is written in compile root or gen-dir dump.
- **L3 e2e**: needs product binary so `-v` nested compile paths appear on captured stdout/stderr
  (in-process CLI buffer does not mirror full nested go test streams the same way).

## Steps

1. Enable UseCLI + Bin for subprocess capture.
2. Descendant copies internal fixture and runs doctest with strategy-specific args.

```go
import (
	"path/filepath"
	"testing"

	"github.com/xhd2015/doctest/session"
	"github.com/xhd2015/doctest/libdoc/testbin"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.UseCLI = true
	req.Bin = testbin.Ensure(t, filepath.Join(d.DOCTEST_ROOT, "..", ".."))
	req.Args = append([]string{"test"}, req.Args...)
	return nil
}
```
