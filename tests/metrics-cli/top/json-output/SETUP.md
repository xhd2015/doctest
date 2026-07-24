# Scenario

**Feature**: `top --json` emits machine-readable JSON

```
metrics top --json --n 3 -> stdout is valid JSON including slow-leaf path
```

## Preconditions

- JSON may be an array of objects or an object wrapping rows; must parse as JSON.

## Steps

1. Seed default fixture.
2. Run `metrics top --json --n 3`.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	ensureFixtureProject(t, req)
	name := "2026-07-16-09-00-00-00-topjson1.jsonl"
	writeRunFile(t, req, name, fixtureRunDefault(runStem(name)))
	req.Args = []string{"metrics", "top", "--json", "--n", "3"}
	return nil
}
```
