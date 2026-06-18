# Scenario

**Feature**: parent not in git, child in git — warn and skip

## Preconditions
- Parent has no `.git`; separate `git init` in child dir only.
- Non-child module paths: parent `parent/tools`, child `parent/cli`.
- Doctest tree in child only.

## Steps
1. Create temp project with child-only git.
2. Run `doctest test -v ./...` from parent root.

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
        git:              gitChildOnly,
    })
    req.WorkDir = projDir
    req.Args = []string{"test", "-v", "./..."}
    return nil
}
```
