# Scenario

**Feature**: leaf without assert import does not add assert replace to nested go.mod

```
# no assert import in SETUP/ASSERT
doctest test --gen-dir outside -> nested go.mod without assert replace
```

## Preconditions

- Leaf ASSERT does not import `github.com/xhd2015/doctest/assert`.

## Steps

1. Create public module without assert import.
2. Run `doctest test <tests> --gen-dir <outsideGenDir> -v`.

```go
import (
	"path/filepath"
	"testing"
)

func Setup(t *testing.T, req *Request) error {
	createPublicModuleProject(t, "", defaultPublicAssertGo())
	req.OutsideGenDir = filepath.Join(t.TempDir(), "generated")
	setupModuleEnv(t, req)
	req.Args = []string{"test", testDir, "--gen-dir", req.OutsideGenDir, "-v"}
	return nil
}
```