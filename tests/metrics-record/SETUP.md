# Scenario

**Feature**: record metrics from doctest test and warn when the default suite is slow

```
# default suite wall clock (in-process)
runner.RunTest(dir, Options{MetricsOn, MetricsRoot, Label*})
  -> if MetricsOn: JSONL run under MetricsRoot
  -> run leaves
  -> if ShouldWarnDefaultSuiteSlow: WARNING on stderr (non-fatal)

# pure helpers (no sleep)
ShouldWarnDefaultSuiteSlow(default, total, elapsed, threshold) -> bool
FormatDefaultSuiteSlowWarning() -> message

# flags
runner.ParseTestOptions(args) -> MetricsOn / remain
```

## Preconditions

- Package `github.com/xhd2015/doctest/libdoc/metrics` provides warn helpers and
  schema constants.
- `core.Options.MetricsOn` / `MetricsRoot` and `runner.ParseTestOptions` /
  `runner.RunTest` wire recording + WARNING.
- Leaves never sleep for 3 minutes; warn cases pass synthetic `elapsed`.
- Recording leaves use `t.TempDir()` as MetricsRoot so user cache is never touched.
- **Layer**: L2 in-process only (no `testbin` / product binary exec).

## Steps

1. Leaf Setup sets `req.Op` and branch fields.
2. Root `Run` dispatches pure helpers, flag parse, or suite recording via
   `runner.RunTest`.
3. Leaf Assert checks boolean, message substrings, option bits, or JSONL events.

## Context

- Default suite ⇔ `!LabelAll && len(LabelExprs)==0`.
- Metrics default **off**; `--metrics-on` is opt-in only.
- WARNING never fails the process by itself.
- Leaf-level events preferred when the go-test JSON path is available.
- Recording leaves labeled `heavy` exercise nested prepare/go test and are
  skipped by default discovery unless `--label heavy` or `--label-all`.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	// Leaves set Op and scenario fields; root only documents shared contract.
	return nil
}
```
