# Scenario

**Feature**: show with run-id dumps that specific run

```
metrics show <older-stem> -> old/only-leaf content
```

## Preconditions

- Older and newer runs both present.

## Steps

1. Seed two runs.
2. Show the older stem explicitly.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	ensureFixtureProject(t, req)
	older := "2026-07-01-10-00-00-00-showid01.jsonl"
	newer := "2026-07-16-09-00-00-00-showid02.jsonl"
	writeRunFile(t, req, older, fixtureRunOlder(runStem(older)))
	writeRunFile(t, req, newer, fixtureRunDefault(runStem(newer)))
	req.Args = []string{"metrics", "show", runStem(older)}
	return nil
}
```
