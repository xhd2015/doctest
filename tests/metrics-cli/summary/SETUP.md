# Scenario

**Feature**: `metrics summary` aggregates recent runs

```
runs/*.jsonl
  -> metrics summary [--last N] [--default-only] [--json]
  -> multi-run aggregate (counts / wall times)
```

## Preconditions

- `--last N` selects the N newest files by filename.

## Steps

1. Seed multiple run fixtures.
2. Run summary with flags under test.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	ensureFixtureProject(t, req)
	return nil
}
```
