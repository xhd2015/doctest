# Scenario

**Feature**: parent-internal + assert via unified layout A (Kind B expose + mapping-gen go.mod replace)

```
# internal/ import detected
doctest -> mapping-gen (layout A) + Kind B expose -> go test ./…/suite
```

## Preconditions

- Fixture module imports `example.com/app/internal/greet` in harness `Run()`.
- Gen root is mapping-gen module `testcase` (not classic `.doctest_run_*`).
- **L3 e2e**: needs product binary so `-v` nested paths appear on captured stdout/stderr
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
