# Scenario

**Feature**: prune deletes oldest files when count exceeds retention (30)

```
35 run files -> metrics prune -> 30 remain (newest)
```

## Preconditions

- Files named with increasing UTC timestamps so sort order is clear.

## Steps

1. Create 35 minimal run JSONL files.
2. Run `metrics prune`.
3. Snapshot remaining files via Response.RunFiles.

```go
import (
	"fmt"
	"testing"
)

func Setup(t *testing.T, req *Request) error {
	ensureFixtureProject(t, req)
	req.SnapshotRunFilesAfter = true
	// 35 files: 2026-06-01 .. padded so lexicographic order matches creation index
	for i := 0; i < 35; i++ {
		// day-of-year style uniqueness within June/July
		name := fmt.Sprintf("2026-06-%02d-12-00-00-00-prune%02d.jsonl", (i%28)+1, i)
		// ensure uniqueness when day wraps: embed i in suffix already
		if i >= 28 {
			name = fmt.Sprintf("2026-07-%02d-12-00-00-00-prune%02d.jsonl", (i-28)+1, i)
		}
		writeRunFile(t, req, name, []map[string]any{
			{"type": "run_start", "run_id": runStem(name)},
			{"type": "run_end", "wall_ns": int64(1), "passed": 0, "total": 0, "skipped": 0, "exit_ok": true},
		})
	}
	req.Args = []string{"metrics", "prune"}
	return nil
}
```
