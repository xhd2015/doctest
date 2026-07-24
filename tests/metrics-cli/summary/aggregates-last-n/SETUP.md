# Scenario

**Feature**: summary `--last 2` aggregates the two newest runs

```
three runs on disk; --last 2 focuses on the two newest basenames
```

## Preconditions

- Three fixtures with distinct stems for assertion.

## Steps

1. Write three run files (old, mid, new).
2. Run `metrics summary --last 2`.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	ensureFixtureProject(t, req)
	a := "2026-07-01-10-00-00-00-sumold01.jsonl"
	b := "2026-07-10-08-00-00-00-summid01.jsonl"
	c := "2026-07-16-09-00-00-00-sumnew01.jsonl"
	writeRunFile(t, req, a, fixtureRunOlder(runStem(a)))
	writeRunFile(t, req, b, fixtureRunDefault(runStem(b)))
	writeRunFile(t, req, c, fixtureRunDefault(runStem(c)))
	req.Args = []string{"metrics", "summary", "--last", "2"}
	return nil
}
```
