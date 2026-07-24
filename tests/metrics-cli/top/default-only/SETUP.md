# Scenario

**Feature**: `top --default-only` uses a default-suite run, not a labeled suite

```
older labeled-suite run (only/labeled-suite-leaf 9s) is newest? NO —
newer is default suite; also place labeled-suite as a second file.
When --default-only, ranking must come from default-suite data
(group/slow-leaf), not only/labeled-suite-leaf if that run is excluded.

If latest is labeled suite, --default-only should fall back to latest
default-suite run among remaining files.
```

## Preconditions

- Two runs: older default with slow-leaf; newer labeled-suite with only/labeled-suite-leaf.
- Without --default-only, latest would be labeled suite; with flag, use default.

## Steps

1. Write older default fixture + newer labeled-suite fixture.
2. Run `metrics top --default-only`.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	ensureFixtureProject(t, req)
	defName := "2026-07-10-08-00-00-00-defsuite.jsonl"
	labName := "2026-07-16-12-00-00-00-labsuite.jsonl"
	writeRunFile(t, req, defName, fixtureRunDefault(runStem(defName)))
	writeRunFile(t, req, labName, fixtureRunLabeledSuite(runStem(labName)))
	req.Args = []string{"metrics", "top", "--default-only"}
	return nil
}
```
