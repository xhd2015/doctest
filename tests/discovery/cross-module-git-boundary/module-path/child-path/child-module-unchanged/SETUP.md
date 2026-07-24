# Scenario

**Feature**: child module path regression guard

## Preconditions
- Single git repo wrapping parent and nested module.
- Parent module `parent/tools`, child module `parent/tools/cli` (child path).
- Doctest tree exists only in child module.

## Steps
1. Create temp project with child module path and single git repo.
2. Run `doctest test -v ./...` from project root.

```go
import (
    "testing"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
    projDir := createModuleProject(t, moduleProjectConfig{
        parentModulePath: "parent/tools",
        childDir:         "cli",
        childModulePath:  "parent/tools/cli",
        childTestName:    "child_test",
        git:              gitSingleRepo,
    })
    req.WorkDir = projDir
    req.Args = []string{"test", "-v", "./..."}
    return nil
}
```
