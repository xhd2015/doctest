# Scenario

**Feature**: stderr gen-root display depends on gen-dir source

```
# gen-dir modes
auto gen-dir -> mapping-gen cache under home | explicit gen-dir -> user path under cwd
```

## Preconditions

- Leaves configure `req.GenDir` (empty for auto, `"_gen"` for explicit).

## Steps

1. Inherit root tree creation and stderr capture.

## Context

- Auto mode uses the mapping-gen cache (not under project cwd).
- Explicit mode uses `_gen` relative to project root.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	// Default to auto mapping-gen; explicit-gen-dir leaf overrides to "_gen".
	req.GenDir = ""
	return nil
}
```