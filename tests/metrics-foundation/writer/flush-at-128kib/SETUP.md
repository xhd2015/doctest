# Scenario

**Feature**: cumulative buffer ≥ 128KiB flushes before Close

```
# threshold
Write(small) + Write(pad ≥ 128KiB) -> Stat before Close > 0 -> Close
```

## Preconditions

- `PadBytes` large enough that encoded JSON exceeds 128KiB (use 128*1024).
- Inspect size before Close.

## Steps

1. Write a tiny run_start then a pad event of ~128KiB blob.
2. Stat before Close — file must already be non-empty (threshold flush).
3. Close and confirm final file still readable.

## Context

- Constant under test: `metrics.FlushThreshold == 128 * 1024`.
- Pad event type may be implementation-defined (`pad`); tests only care about flush side effect.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.InspectBeforeClose = true
	req.PadBytes = 128 * 1024
	req.Events = []map[string]any{
		{
			"type":       "run_start",
			"run_id":     "run-flush-1",
			"ts":         "2026-07-17T02:00:00Z",
			"project_id": "writer_proj",
			"cwd":        "/w",
			"argv":       []any{"doctest"},
		},
	}
	return nil
}
```
