# Scenario

**Feature**: lifelog layout mirror discovers nested module tests

## Preconditions
- Root `go.mod`: `module github.com/xhd2015/lifelog/tools`.
- `lifelog-cli/go.mod`: `module github.com/xhd2015/lifelog/lifelog-cli`.
- Doctest at `lifelog-cli/tests/skill-cli/` only.
- Single git repo wrapping the whole tree.

## Steps
1. Create lifelog-mirror temp project.
2. Run `doctest test -v ./...` from mirrored root.

```go
import (
    "testing"
)

func Setup(t *testing.T, req *Request) error {
    projDir := createLifelogMirrorProject(t)
    req.WorkDir = projDir
    req.Args = []string{"test", "-v", "./..."}
    return nil
}
```
