# Scenario

**Feature**: single git repo discovers sibling module (lifelog reproduction)

## Preconditions
- Single git repo initialized at project root.
- Parent `module parent/tools`, child `module parent/cli` (non-child).
- Doctest tree in child only.

## Steps
1. Create temp project with single git repo wrapping both modules.
2. Run `doctest test -v ./...` from project root.

```go
import (
    "testing"
)

func Setup(t *testing.T, req *Request) error {
    projDir := createModuleProject(t, moduleProjectConfig{
        parentModulePath: "parent/tools",
        childDir:         "cli",
        childModulePath:  "parent/cli",
        childTestName:    "child_test",
        git:              gitSingleRepo,
    })
    req.WorkDir = projDir
    req.Args = []string{"test", "-v", "./..."}
    return nil
}
```
