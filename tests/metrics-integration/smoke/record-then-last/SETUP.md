# Scenario

**Feature**: after a tiny suite run, `metrics last` summarizes the recorded run

```
RunTest(1-leaf, MetricsRoot=tmp)
  -> JSONL under MetricsRoot
  -> doctest metrics last (DOCTEST_METRICS_ROOT, cwd=fixture)
  -> exit 0 + run summary evidence
```

## Preconditions

- Empty MetricsRoot before record.
- Analyze args: `metrics last`.

## Steps

1. Record one-leaf suite with metrics on.
2. Run `metrics last`.
3. Assert new JSONL and analyze exit 0 with non-empty useful output.

## Context

- Output may mention run id/filename stem, pass/total counts, or leaf path.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.AnalyzeArgs = []string{"metrics", "last"}
	return nil
}
```
