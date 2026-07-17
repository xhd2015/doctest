# Scenario

**Feature**: small events stay buffered until Close

```
# sub-threshold
Write(run_start) -> Stat before Close ≈ 0 -> Close -> file has JSONL
```

## Preconditions

- A single small run_start event (well under 128KiB).
- Inspect file size before Close.

## Steps

1. Write one small event.
2. Stat path before Close (expect empty / zero size / missing).
3. Close and confirm data is on disk.

## Context

- Flush threshold is 128KiB; one short JSON line must not force mid flush.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.InspectBeforeClose = true
	req.Events = []map[string]any{
		{
			"type":       "run_start",
			"run_id":     "run-small-1",
			"ts":         "2026-07-17T01:00:00Z",
			"project_id": "writer_proj",
			"cwd":        "/w",
			"argv":       []any{"doctest"},
		},
	}
	return nil
}
```
