## Preconditions
- A valid doctest tree exists.
- The root Setup has built the doctest binary.

## Steps
1. Run `doctest test <dir> --gen-dir <tmp>`.
2. Verify generated files exist at the specified directory.

```go
import (
    "path/filepath"
    "testing"
)

func Setup(t *testing.T, req *Request) error {
    testDir := createTempTestProject(t, "mytest")
    genDir := filepath.Join(t.TempDir(), "generated")
    req.Args = []string{"test", testDir, "--gen-dir", genDir}
    req.Env = append(req.Env, "GEN_DIR="+genDir)
    return nil
}
```
