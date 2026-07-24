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
	"testing"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	createPublicModuleProject(t, req, "", defaultSessionAssertGo(), true)
	setupModuleEnv(t, req)
	// Request-local gen dir: Parallel-safe (no package-level var genDir).
	req.GenDir = t.TempDir()
	req.Args = []string{"test", req.TestDir, "--gen-dir", req.GenDir, "-v"}
	return nil
}
```
