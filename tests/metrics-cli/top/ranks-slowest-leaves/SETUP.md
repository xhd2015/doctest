# Scenario

**Feature**: top ranks leaves slowest-first from latest run

```
default fixture leaves (5s, 3s labeled, 2s, 0.1s)
  -> metrics top
  -> group/slow-leaf appears before mid-leaf and fast-leaf
```

## Preconditions

- Single latest default-suite run fixture.

## Steps

1. Write default fixture as newest run.
2. Run `metrics top` (defaults).

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	ensureFixtureProject(t, req)
	name := "2026-07-16-09-00-00-00-topbase1.jsonl"
	writeRunFile(t, req, name, fixtureRunDefault(runStem(name)))
	req.Args = []string{"metrics", "top"}
	return nil
}
```
