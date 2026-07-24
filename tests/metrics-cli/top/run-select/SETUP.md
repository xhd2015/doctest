# Scenario

**Feature**: `top --run <id>` selects a non-latest run file

```
older run (old/only-leaf) + newer default run
  -> metrics top --run <older-stem>
  -> ranks old/only-leaf, not group/slow-leaf
```

## Preconditions

- Run id is the JSONL basename without `.jsonl`.

## Steps

1. Write two run files.
2. Run top with `--run` set to the older stem.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	ensureFixtureProject(t, req)
	older := "2026-07-01-10-00-00-00-pickold1.jsonl"
	newer := "2026-07-16-09-00-00-00-picknew1.jsonl"
	writeRunFile(t, req, older, fixtureRunOlder(runStem(older)))
	writeRunFile(t, req, newer, fixtureRunDefault(runStem(newer)))
	req.Args = []string{"metrics", "top", "--run", runStem(older)}
	return nil
}
```
