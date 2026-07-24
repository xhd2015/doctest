# Scenario

**Feature**: a doctest tree exists with a stub Run() that returns "error not implemented"

```
# full TDD cycle: design -> RED -> seal -> implement -> GREEN
orchestrator -> design agent -> writes tests -> RED (all fail)

# seal tests, hand off to implementer
orchestrator -> git add tests/ -> implement agent -> writes code -> GREEN (all pass)

# question/answer loop
user <--questions-- implement agent <--yields-- orchestrator -> resume
```

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

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.UseCLI = true // true e2e nested doctest test
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
