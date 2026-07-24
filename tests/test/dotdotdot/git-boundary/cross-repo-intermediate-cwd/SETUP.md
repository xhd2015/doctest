# Scenario

**Feature**: parent dotdotdot helpers (createTestTree etc.) are available

```
# build and run test binary, report results
doctest test <dir> -> build -> run binary -> pass/fail per leaf -> exit code

# path patterns
.../ -> walk tree | subdir -> run subtree | multi-dir -> aggregate results

# output
progress dots -> . F | verbose -> go test -v | count -> N tests
```

## Preconditions
- Parent dotdotdot helpers (createTestTree etc.) are available.
- The CWD is a directory inside a nested git repo between a parent Go module and a child Go module with tests.

## Steps
1. Create a parent git repo with go.mod (no doctests).
2. Inside parent, create `ext/` — a different git repo (no go.mod, no tests).
3. Inside ext, create `inner/` — with go.mod and a DOCTEST.md tree.
4. Run `doctest test -v ./...` from ext/.
5. Verify the inner tests are discovered (walk up stops at ext git boundary, then falls through to subdir discovery).

```go
import (
    "os"
    "os/exec"
    "path/filepath"
    "testing"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
    parent := t.TempDir()
    if err := initGitRepo(parent); err != nil {
        t.Fatalf("init parent repo: %v", err)
    }
    if err := writeGoMod(parent, "testproj"); err != nil {
        t.Fatalf("write parent go.mod: %v", err)
    }

    extDir := filepath.Join(parent, "ext")
    if err := os.MkdirAll(extDir, 0755); err != nil {
        t.Fatalf("mkdir ext: %v", err)
    }
    if err := initGitRepo(extDir); err != nil {
        t.Fatalf("init ext repo: %v", err)
    }

    innerDir := filepath.Join(extDir, "inner")
    if err := os.MkdirAll(innerDir, 0755); err != nil {
        t.Fatalf("mkdir inner: %v", err)
    }
    if err := writeGoMod(innerDir, "inner"); err != nil {
        t.Fatalf("write inner go.mod: %v", err)
    }
    if err := createTestTree(innerDir, "inner_test"); err != nil {
        t.Fatalf("create inner_test: %v", err)
    }

    req.WorkDir = extDir
    req.Args = []string{"test", "-v", "./..."}
    return nil
}
```
