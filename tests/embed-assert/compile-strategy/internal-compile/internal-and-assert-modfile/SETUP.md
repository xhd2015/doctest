# Scenario

**Feature**: leaf importing both internal and assert compiles via -modfile with assert replace

```
# internal + assert imports
doctest test -v -> .doctest_run_* -> go test -modfile=<tmp> with assert replace
```

## Preconditions

- Fixture `internal-assert-module` imports internal greet and assert in ASSERT.

## Steps

1. Copy internal+assert fixture.
2. Run `doctest test <tests> -v`.

```go
import (
"testing"
"github.com/xhd2015/doctest/session"
)
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	createInternalAssertProject(t, d, req)
	setupModuleEnv(t, req)
	req.Args = []string{"test", req.TestDir, "-v"}
	return nil
}
```