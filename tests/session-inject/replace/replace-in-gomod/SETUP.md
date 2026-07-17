# Scenario

**Feature**: generated nested go.mod contains session replace pointing at session-mod cache

```
doctest test --gen-dir <dir> with session import
  -> testcase go.mod has replace session => cache
```

## Preconditions

- Leaf imports session package.
- `--gen-dir` captures generated module for inspection.

## Steps

1. Create public module with session import.
2. Run `doctest test <tests> --gen-dir <tmp> -v`.
3. Find nested go.mod and assert replace line.

```go
import (
	"path/filepath"
	"testing"
)

func Setup(t *testing.T, req *Request) error {
	createPublicModuleProject(t, "", defaultSessionAssertGo(), true)
	setupModuleEnv(t, req)
	genDir = t.TempDir()
	req.Args = []string{"test", testDir, "--gen-dir", genDir, "-v"}
	return nil
}
```
