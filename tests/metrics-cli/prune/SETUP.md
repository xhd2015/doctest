# Scenario

**Feature**: `metrics prune` enforces run-file retention (keep newest 30)

```
runs/*.jsonl count C
  -> metrics prune
  -> keep min(C, 30) newest files; delete older excess
```

## Preconditions

- Retention constant `DefaultRunRetention` = 30 (see root DOCTEST.md).
- Ordering by filename lexicographic (UTC names).

## Steps

1. Seed N empty or minimal run files.
2. Run prune; snapshot remaining basenames.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	ensureFixtureProject(t, req)
	req.SnapshotRunFilesAfter = true
	return nil
}
```
