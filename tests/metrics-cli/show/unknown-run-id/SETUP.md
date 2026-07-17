# Scenario

**Feature**: show with unknown run-id exits non-zero

```
metrics show definitely-missing-run-id-xyz -> non-zero
```

## Preconditions

- Optional: one real run present so store exists; id still missing.

## Steps

1. Seed one real run (so missing-id is not conflated with empty store).
2. Show a nonexistent id.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	ensureFixtureProject(t, req)
	name := "2026-07-16-09-00-00-00-showmiss.jsonl"
	writeRunFile(t, req, name, fixtureRunDefault(runStem(name)))
	req.Args = []string{"metrics", "show", "definitely-missing-run-id-xyz"}
	return nil
}
```
