# Scenario

**Feature**: skipped leaves may be recorded as leaf_end only (no leaf_start)

```
# skip without start
run_start -> leaf_end(result=skip, optional ts_start omitted) -> run_end
```

## Preconditions

- Sequence has no `leaf_start` for the skipped path.

## Steps

1. Write run_start, a skip leaf_end without ts_start, and run_end.
2. Decode and verify skip contract.

## Context

- `elapsed_ns` may be 0 for pure skips.
- Omitting `ts_start` is allowed for skip-only ends.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Events = []map[string]any{
		{
			"type":       "run_start",
			"run_id":     "run-skip-1",
			"ts":         "2026-07-17T12:00:00Z",
			"project_id": "github.com_xhd2015_doctest",
			"cwd":        "/work",
			"argv":       []any{"doctest", "test"},
		},
		{
			"type":       "leaf_end",
			"path":       "slow/ui-automation",
			"ts_end":     "2026-07-17T12:00:00.001Z",
			"elapsed_ns": int64(0),
			"result":     "skip",
			"cached":     false,
			"labels":     []any{"heavy"},
		},
		{
			"type":    "run_end",
			"wall_ns": int64(1_000_000),
			"passed":  0,
			"total":   1,
			"skipped": 1,
			"exit_ok": true,
		},
	}
	return nil
}
```
