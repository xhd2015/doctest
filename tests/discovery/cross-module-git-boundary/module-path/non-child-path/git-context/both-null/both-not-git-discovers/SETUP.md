# Scenario

**Feature**: both-not-git discovers sibling module

## Preconditions
- Parent: no git, `go.mod` `module parent/tools`.
- Child: no git, `go.mod` `module parent/cli` (non-child).
- Doctest tree in child only.

## Steps
1. Create temp project with no git repos.
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
        git:              gitNone,
    })
    req.WorkDir = projDir
    req.Args = []string{"test", "-v", "./..."}
    return nil
}
```
