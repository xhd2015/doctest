# Scenario

**Feature**: `metrics top` ranks leaf_end events by elapsed_ns

```
runs/*.jsonl leaf_end{path,elapsed_ns,labels*}
  -> metrics top [--n N] [--unlabeled-only] [--default-only] [--json] [--run last|ID]
  -> ranked slow leaves
```

## Preconditions

- Ranking uses `elapsed_ns` descending from the selected run.
- Labels for `--unlabeled-only` come from `leaf_start.labels` for the path.

## Steps

1. Seed fixtures for the leaf scenario.
2. Run `metrics top` with flags under test.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	ensureFixtureProject(t, req)
	return nil
}
```
