# Scenario

**Feature**: internal-only leaf (no author assert) still uses -modfile for always-on assert+session

```
# internal import only, no author assert import
# always-on assertImport + sessionImport → WriteInternalModfile + -modfile=
doctest test -v -> temp -modfile (parent go.mod + assert/session replaces)
```

## Preconditions

- Fixture `internal-only-module` has no assert import.
- Always-on assert/session still wires a temp modfile for internal-compile.

## Steps

1. Copy internal-only fixture.
2. Run `doctest test <tests> -v`.

```go
import (
"testing"
"github.com/xhd2015/doctest/session"
)
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	createInternalOnlyProject(t, d)
	setupModuleEnv(t, req)
	req.Args = []string{"test", testDir, "-v"}
	return nil
}
```
