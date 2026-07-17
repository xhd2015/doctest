# Scenario

**Feature**: `top --n` limits the number of ranked leaves

```
metrics top --n 2 -> at most two leaf paths from the ranking
```

## Preconditions

- Fixture has ≥ 3 leaf_end events.

## Steps

1. Seed default fixture.
2. Run `metrics top --n 2`.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	ensureFixtureProject(t, req)
	name := "2026-07-16-09-00-00-00-topn0002.jsonl"
	writeRunFile(t, req, name, fixtureRunDefault(runStem(name)))
	req.Args = []string{"metrics", "top", "--n", "2"}
	return nil
}
```
