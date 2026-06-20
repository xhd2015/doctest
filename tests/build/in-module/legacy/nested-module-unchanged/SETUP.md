# Scenario

**Feature**: public import with gen-dir outside module uses legacy nested testcase module

```
# legacy: no internal import, gen-dir outside -> nested go.mod (testcase + replace)
doctest test --gen-dir <outside>/generated -> module testcase -> PASS
```

## Preconditions

- Temp module with public `pkg/greet` (no internal import).
- `--gen-dir` is outside the parent module.

## Steps

1. Create public-import temp module.
2. Run `doctest test <tests> --gen-dir <outsideGenDir> -v`.

```go
import (
	"path/filepath"
	"testing"
)

var outsideGenDir string

func Setup(t *testing.T, req *Request) error {
	createPublicModuleProject(t)
	_ = genDir
	outsideGenDir = filepath.Join(t.TempDir(), "generated")
	req.WorkDir = moduleRoot
	req.Args = []string{"test", testDir, "--gen-dir", outsideGenDir, "-v"}
	return nil
}
```