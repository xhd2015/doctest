# Scenario

**Feature**: `metrics show` dumps one run file

```
metrics show [run-id]
  omit id -> latest
  with id -> that runs/*.jsonl
```

## Preconditions

- Run id matches basename without `.jsonl`.

## Steps

1. Seed one or more run fixtures.
2. Run show with optional id.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	ensureFixtureProject(t, req)
	return nil
}
```
