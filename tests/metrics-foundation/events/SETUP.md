# Scenario

**Feature**: schema_version 1 JSONL event types encode and decode

```
# event types
run_start -> leaf_start -> leaf_end(pass|fail|skip) -> run_end

# skip path
skip may be leaf_end only (no leaf_start)

# git metadata
run_start may include branch/commit; must not include dirty flag
```

## Preconditions

- Writer creates a run file under a temp cache.
- Events are maps/objects with `type` and `schema_version`.

## Steps

1. Leaf builds an event sequence.
2. Run writes via `OpenWriter` / `Write` / `Close`.
3. Assert re-reads JSON lines and field contracts.

## Context

- Every line is one JSON object ending with `\n`.
- Field names use snake_case as in the requirement.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Op = "write_sequence"
	req.CacheDir = t.TempDir()
	req.ProjectID = "github.com_xhd2015_doctest"
	req.Suffix = "evt00001"
	return nil
}
```
