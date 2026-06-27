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
import "testing"

func Setup(t *testing.T, req *Request) error {
	createInternalAssertProject(t)
	setupModuleEnv(t, req)
	req.Args = []string{"test", testDir, "-v"}
	return nil
}
```