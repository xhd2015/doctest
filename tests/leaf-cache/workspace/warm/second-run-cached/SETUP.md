# Scenario

**Feature**: second default workspace run reports leaf-cache Cached > 0

```
run1: workspace two trees -> exit 0, may be 0 Cached (cold)
run2: workspace two trees -> exit 0, Cached >= 2 (both leaves warm)
```

## Preconditions

- Parent prepared all-pass multi-tree fixture and identical Args/Args2.

## Steps

1. Keep double-run configuration.
2. Assert run2 Cached >= 2 (one skip per tree leaf) and both exits 0.

## Context

- Proves multi-prep prepare skip env + RecordPasses + summary Cached on
  `finishWorkspaceGoTest` (not go testcache; fresh GOCACHE per run).

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	t.Helper()
	req.Op = "runtime_multi"
	return nil
}
```
