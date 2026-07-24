# Scenario

**Feature**: summary `--default-only` ignores non-default suite runs

```
default run + labeled-suite run
  -> metrics summary --last 10 --default-only
  -> aggregate only default-suite file(s)
```

## Preconditions

- Labeled suite run is present but excluded by the flag.

## Steps

1. Write one default + one labeled-suite fixture.
2. Run `metrics summary --default-only --last 5`.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	ensureFixtureProject(t, req)
	defName := "2026-07-10-08-00-00-00-sumdef01.jsonl"
	labName := "2026-07-16-12-00-00-00-sumlab01.jsonl"
	writeRunFile(t, req, defName, fixtureRunDefault(runStem(defName)))
	writeRunFile(t, req, labName, fixtureRunLabeledSuite(runStem(labName)))
	req.Args = []string{"metrics", "summary", "--default-only", "--last", "5"}
	return nil
}
```
