# Scenario

**Feature**: summary `--json` emits valid JSON

```
metrics summary --last 1 --json -> parseable JSON on stdout
```

## Preconditions

- At least one run fixture.

## Steps

1. Seed one default run.
2. Run `metrics summary --last 1 --json`.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	ensureFixtureProject(t, req)
	name := "2026-07-16-09-00-00-00-sumjson1.jsonl"
	writeRunFile(t, req, name, fixtureRunDefault(runStem(name)))
	req.Args = []string{"metrics", "summary", "--last", "1", "--json"}
	return nil
}
```
