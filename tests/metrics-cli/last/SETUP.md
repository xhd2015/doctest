# Scenario

**Feature**: `metrics last` summarizes the newest run file

```
runs/*.jsonl (lexicographic UTC names)
  -> doctest metrics last
  -> newest run summary (id + counts)
```

## Preconditions

- Newest file is the maximum basename under `runs/`.

## Steps

1. Seed WorkDir + MetricsRoot.
2. Leaf writes zero or more run fixtures.
3. Run `metrics last`.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	ensureFixtureProject(t, req)
	return nil
}
```
