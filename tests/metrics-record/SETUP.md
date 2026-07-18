# Scenario

**Feature**: record metrics from doctest test and warn when the default suite is slow

```
# default suite wall clock
doctest test <dir> -> Options(MetricsOn, MetricsRoot, Label*)
  -> if MetricsOn: JSONL run under MetricsRoot
  -> run leaves
  -> if ShouldWarnDefaultSuiteSlow: WARNING on stderr (non-fatal)

# pure helpers (no sleep)
ShouldWarnDefaultSuiteSlow(default, total, elapsed, threshold) -> bool
FormatDefaultSuiteSlowWarning() -> message
```

## Preconditions

- P1 package `github.com/xhd2015/doctest/libdoc/metrics` exists (writer/paths/id).
- P2 symbols expected (implementer contract; tree is RED until present):
  - `metrics.DefaultSuiteWarnThreshold` (`3 * time.Minute`)
  - `metrics.ShouldWarnDefaultSuiteSlow(defaultSuite bool, total int, elapsed, threshold time.Duration) bool`
  - `metrics.FormatDefaultSuiteSlowWarning() string`
  - `core.Options.MetricsOn bool` and `core.Options.MetricsRoot string`
  - `runner.ParseTestOptions(args []string) (core.Options, []string, error)` understands `--metrics-on`
  - `runner.RunTest(dir string, opts core.Options) error` runs one tree with metrics + WARNING wiring
    (CLI path may also honor env `DOCTEST_METRICS_ROOT` as MetricsRoot override)
- Leaves never sleep for 3 minutes; warn cases pass synthetic `elapsed`.
- Recording leaves use `t.TempDir()` as MetricsRoot so user cache is never touched.

## Steps

1. Leaf Setup sets `req.Op` and branch fields.
2. Root `Run` dispatches pure helpers, flag parse, or suite recording.
3. Leaf Assert checks boolean, message substrings, option bits, or JSONL events.

## Context

- Default suite ⇔ `!LabelAll && len(LabelExprs)==0`.
- Metrics default **off**; `--metrics-on` is opt-in only.
- WARNING never fails the process by itself.
- Leaf-level events preferred when the go-test JSON path is available.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	// Leaves set Op and scenario fields; root only documents shared contract.
	return nil
}
```
