# Scenario

**Feature**: after a tiny suite run, `metrics top` ranks leaves from the recorded run

```
RunTest(1-leaf, MetricsRoot=tmp)
  -> JSONL with leaf_end for a_pass_0
  -> doctest metrics top --n 5
  -> exit 0 + leaf path evidence
```

## Preconditions

- Empty MetricsRoot before record.
- Analyze args: `metrics top --n 5`.

## Steps

1. Record one-leaf suite with metrics on.
2. Run `metrics top --n 5`.
3. Assert JSONL and top output references the fixture leaf path.

## Context

- Fixture leaf name from `WritePassFailTree` is `a_pass_0`.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.AnalyzeArgs = []string{"metrics", "top", "--n", "5"}
	return nil
}
```
