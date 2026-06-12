## Preconditions
- A doctest tree exists with a stub Run() that returns "error not implemented".

## Steps
1. Create a doctest tree with stub Run().
2. Run `doctest test -v` on it.
3. Expect all tests to fail (RED).

```go
import (
    "os"
    "os/exec"
    "path/filepath"
    "testing"
)

func Setup(t *testing.T, req *Request) error {
    dir := filepath.Join(t.TempDir(), "test-tree")
    createDoctestTree(t, dir, true) // stub = true
    req.Env = append(req.Env, "TEST_TREE_DIR="+dir)

    doctestBin := ""
    for _, env := range req.Env {
        if len(env) > 12 && env[:12] == "DOCTEST_BIN=" {
            doctestBin = env[12:]
            break
        }
    }
    if doctestBin == "" {
        t.Fatal("DOCTEST_BIN not set by parent")
    }

    req.Args = []string{"test", "-v", dir}
    _ = exec.Command // ensure import
    return nil
}
```
