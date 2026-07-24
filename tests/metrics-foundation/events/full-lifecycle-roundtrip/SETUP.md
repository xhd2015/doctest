# Scenario

**Feature**: full run lifecycle encodes to multi-line JSONL and decodes cleanly

```
# sequence
run_start
leaf_start (path A)
leaf_end pass (path A)
leaf_end skip (path B)   # skip without prior start is allowed in a later leaf; here we still write end only for B after A
run_end
```

## Preconditions

- Five events covering all four types (leaf_end used twice: pass + skip).

## Steps

1. Write run_start, leaf_start, leaf_end(pass), leaf_end(skip), run_end.
2. Close writer and read the file.
3. Decode each line as JSON.

## Context

- `schema_version` is 1 on every event.
- Trailing newline per line (JSONL).

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Events = []map[string]any{
		{
			"type":       "run_start",
			"run_id":     "run-lifecycle-1",
			"ts":         "2026-07-17T12:34:56Z",
			"project_id": "github.com_xhd2015_doctest",
			"cwd":        "/work/repo",
			"argv":       []any{"doctest", "test", "./..."},
			"git_branch": "master",
			"git_commit": "abc1234",
			"session_id": "sess-1",
		},
		{
			"type":     "leaf_start",
			"path":     "project-id/https-origin",
			"root":     "tests/metrics-foundation",
			"ts_start": "2026-07-17T12:34:56.001Z",
			"labels":   []any{},
		},
		{
			"type":       "leaf_end",
			"path":       "project-id/https-origin",
			"ts_start":   "2026-07-17T12:34:56.001Z",
			"ts_end":     "2026-07-17T12:34:56.010Z",
			"elapsed_ns": int64(9_000_000),
			"result":     "pass",
			"cached":     false,
		},
		{
			"type":       "leaf_end",
			"path":       "project-id/skipped-leaf",
			"ts_end":     "2026-07-17T12:34:56.020Z",
			"elapsed_ns": int64(0),
			"result":     "skip",
			"cached":     false,
		},
		{
			"type":     "run_end",
			"wall_ns":  int64(50_000_000),
			"passed":   1,
			"total":    2,
			"skipped":  1,
			"exit_ok":  true,
			"warnings": []any{},
		},
	}
	return nil
}
```
