# Scenario

**Feature**: leaf without author assert import still gets assert replace in nested go.mod

```
# no assert import in SETUP/ASSERT author code
# always-on assertImport + sessionImport for external modules
doctest test --gen-dir outside -> nested go.mod with assert (+ session) replace
```

## Preconditions

- Leaf ASSERT does not import `github.com/xhd2015/doctest/assert`.
- Outside gen-dir triggers nested testcase module generation (`WriteGoMod`).

## Steps

1. Create public module without assert import.
2. Run `doctest test <tests> --gen-dir <outsideGenDir> -v`.

```go
import (
	"path/filepath"
	"testing"
)

func Setup(t *testing.T, req *Request) error {
	createPublicModuleProject(t, req, "", defaultPublicAssertGo())
	req.OutsideGenDir = filepath.Join(t.TempDir(), "generated")
	setupModuleEnv(t, req)
	req.Args = []string{"test", req.TestDir, "--gen-dir", req.OutsideGenDir, "-v"}
	return nil
}
```
