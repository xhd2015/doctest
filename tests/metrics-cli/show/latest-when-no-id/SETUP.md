# Scenario

**Feature**: show without id dumps the latest run

```
older + newer fixtures -> metrics show -> newer content (group/slow-leaf)
```

## Preconditions

- Two run files present.

## Steps

1. Write older + newer fixtures.
2. Run `metrics show` with no run-id argument.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	ensureFixtureProject(t, req)
	older := "2026-07-01-10-00-00-00-showold1.jsonl"
	newer := "2026-07-16-09-00-00-00-shownew1.jsonl"
	writeRunFile(t, req, older, fixtureRunOlder(runStem(older)))
	writeRunFile(t, req, newer, fixtureRunDefault(runStem(newer)))
	req.Args = []string{"metrics", "show"}
	return nil
}
```
