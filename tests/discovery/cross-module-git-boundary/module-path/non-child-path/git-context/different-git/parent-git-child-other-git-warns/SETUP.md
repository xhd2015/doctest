# Scenario

**Feature**: separate git repos warn and skip sibling module

## Preconditions
- Parent git repo at root, separate `git init` in child dir.
- Non-child module paths: parent `parent/tools`, child `parent/cli`.
- Doctest tree in child only.

## Steps
1. Create temp project with separate git repos.
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
        git:              gitSeparateRepos,
    })
    req.WorkDir = projDir
    req.Args = []string{"test", "-v", "./..."}
    return nil
}
```
