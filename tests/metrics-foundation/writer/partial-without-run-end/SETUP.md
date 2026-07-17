# Scenario

**Feature**: file without run_end remains a readable partial JSONL prefix

```
# incomplete run
run_start + leaf_start + leaf_end(pass)  # no run_end
  -> Close flush -> three decodable lines
```

## Preconditions

- Sequence intentionally omits `run_end`.

## Steps

1. Write three events without run_end.
2. Close writer.
3. Decode all lines; assert types and that none is run_end.

## Context

- Partial metrics files are valid inputs for later analysis tools.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Events = []map[string]any{
		{
			"type":       "run_start",
			"run_id":     "run-partial-1",
			"ts":         "2026-07-17T03:00:00Z",
			"project_id": "writer_proj",
			"cwd":        "/w",
			"argv":       []any{"doctest", "test"},
		},
		{
			"type":     "leaf_start",
			"path":     "a/b",
			"root":     "tests",
			"ts_start": "2026-07-17T03:00:00.001Z",
		},
		{
			"type":       "leaf_end",
			"path":       "a/b",
			"ts_start":   "2026-07-17T03:00:00.001Z",
			"ts_end":     "2026-07-17T03:00:00.002Z",
			"elapsed_ns": int64(1_000_000),
			"result":     "pass",
			"cached":     false,
		},
	}
	return nil
}
```
