# Scenario

**Feature**: no test tree exists yet

```
# full TDD cycle: design -> RED -> seal -> implement -> GREEN
orchestrator -> design agent -> writes tests -> RED (all fail)

# seal tests, hand off to implementer
orchestrator -> git add tests/ -> implement agent -> writes code -> GREEN (all pass)

# question/answer loop
user <--questions-- implement agent <--yields-- orchestrator -> resume
```

## Preconditions
- No test tree exists yet.

## Steps
1. Call createDoctestTree to write the test tree to a temp dir.
2. Pass the dir path to ASSERT via req.Env.

```go
import (
    "os"
    "path/filepath"
    "testing"
)

func Setup(t *testing.T, req *Request) error {
    dir := filepath.Join(t.TempDir(), "test-tree")
    createDoctestTree(t, dir, false)
    req.Env = append(req.Env, "TEST_TREE_DIR="+dir)
    req.Args = []string{"--help"}
    return nil
}
```
