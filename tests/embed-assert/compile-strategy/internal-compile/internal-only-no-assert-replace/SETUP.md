# Scenario

**Feature**: internal-only leaf does not add assert replace to temp modfile

```
# internal import only, no assert
doctest test -v -> -modfile without assert replace
```

## Preconditions

- Fixture `internal-only-module` has no assert import.

## Steps

1. Copy internal-only fixture.
2. Run `doctest test <tests> -v`.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	createInternalOnlyProject(t)
	setupModuleEnv(t, req)
	req.Args = []string{"test", testDir, "-v"}
	return nil
}
```