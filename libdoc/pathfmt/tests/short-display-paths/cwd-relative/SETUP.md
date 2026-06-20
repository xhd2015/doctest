# Scenario

**Feature**: paths under the cwd subtree display as `"."` or `"./rel"`

```
# cwd rules (checked first)
path == cwd -> "." | strict child of cwd -> "./" + rel
```

## Preconditions

- A temp project root exists; process cwd is changed to that root before each leaf runs.

## Steps

1. Save and restore the original cwd.
2. Create a temp project directory with subdirectories as needed per leaf.
3. `chdir` to the project root (or nested cwd for `parent-path`).

## Context

- `req.Path` is set by each leaf to the absolute or relative path under test.

```go
import (
	"path/filepath"
	"testing"
)

func Setup(t *testing.T, req *Request) error {
	saveAndRestoreCwd(t)
	projRoot := t.TempDir()
	mkdirAll(t, filepath.Join(projRoot, "child"))
	mkdirAll(t, filepath.Join(projRoot, "a", "b", "c"))
	req.Path = projRoot
	chdirTo(t, projRoot)
	return nil
}
```