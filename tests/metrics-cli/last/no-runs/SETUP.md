# Scenario

**Feature**: last fails clearly when no run files exist

```
empty metrics project (or missing runs/)
  -> metrics last
  -> non-zero + no-runs message
```

## Preconditions

- Project identity resolves, but `runs/` has no `*.jsonl`.

## Steps

1. Create empty runs directory (or leave missing).
2. Run `metrics last`.

```go
import (
	"os"
	"testing"
)

func Setup(t *testing.T, req *Request) error {
	ensureFixtureProject(t, req)
	// Empty runs dir documents "project known, no files"
	if err := os.MkdirAll(projectRunsDir(req), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	req.Args = []string{"metrics", "last"}
	return nil
}
```
