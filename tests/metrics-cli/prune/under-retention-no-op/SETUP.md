# Scenario

**Feature**: prune is a no-op when run count is under retention

```
3 run files -> metrics prune -> still 3 files
```

## Preconditions

- Count < DefaultRunRetention (30).

## Steps

1. Write 3 minimal run files.
2. Run prune; expect count unchanged.

```go
import (
	"fmt"
	"testing"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	ensureFixtureProject(t, req)
	req.SnapshotRunFilesAfter = true
	for i := 0; i < 3; i++ {
		name := fmt.Sprintf("2026-07-0%d-12-00-00-00-keep%02d.jsonl", i+1, i)
		writeRunFile(t, req, name, []map[string]any{
			{"type": "run_start", "run_id": runStem(name)},
			{"type": "run_end", "wall_ns": int64(1), "passed": 0, "total": 0, "skipped": 0, "exit_ok": true},
		})
	}
	req.Args = []string{"metrics", "prune"}
	return nil
}
```
