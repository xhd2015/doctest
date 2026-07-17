# Scenario

**Feature**: last picks the newest of multiple run files

```
older run + newer default run
  -> metrics last
  -> summary of newer (not old/only-leaf)
```

## Preconditions

- Two run files; newer filename sorts after older.

## Steps

1. Write older fixture and newer default-suite fixture.
2. Run `metrics last`.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	ensureFixtureProject(t, req)
	older := "2026-07-01-10-00-00-00-oldrun01.jsonl"
	newer := "2026-07-16-09-00-00-00-newrun01.jsonl"
	writeRunFile(t, req, older, fixtureRunOlder(runStem(older)))
	writeRunFile(t, req, newer, fixtureRunDefault(runStem(newer)))
	req.Args = []string{"metrics", "last"}
	return nil
}
```
