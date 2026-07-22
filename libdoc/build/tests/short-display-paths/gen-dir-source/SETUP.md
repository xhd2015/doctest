# Scenario

**Feature**: stderr gen-root display depends on gen-dir source (no process Chdir)

```
# gen-dir modes (Parallel-safe)
auto gen-dir -> mapping-gen cache under home (~/ display)
explicit gen-dir -> absolute path under sandbox project (Short may be abs)
# never: os.Chdir(project) to force relative _gen / tests/feature
```

## Preconditions

- Leaves configure `req.GenDir` (empty for auto, `"_gen"` for explicit).
- Root `Run` resolves relative GenDir under the sandbox project.

## Steps

1. Inherit root tree creation and stderr capture without Chdir.

## Context

- Auto mode uses the mapping-gen cache (not under project cwd).
- Explicit mode uses `_gen` joined under sandbox `projRoot` as absolute.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	// Default to auto mapping-gen; explicit-gen-dir leaf overrides to "_gen".
	req.GenDir = ""
	return nil
}
```
