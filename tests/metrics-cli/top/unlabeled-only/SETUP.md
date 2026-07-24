# Scenario

**Feature**: `top --unlabeled-only` drops leaves that have labels

```
metrics top --unlabeled-only
  -> includes group/slow-leaf
  -> excludes group/labeled-leaf (labels: slow)
```

## Preconditions

- Fixture includes labeled and unlabeled leaves.

## Steps

1. Seed default fixture.
2. Run `metrics top --unlabeled-only`.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	ensureFixtureProject(t, req)
	name := "2026-07-16-09-00-00-00-topunlab.jsonl"
	writeRunFile(t, req, name, fixtureRunDefault(runStem(name)))
	req.Args = []string{"metrics", "top", "--unlabeled-only"}
	return nil
}
```
