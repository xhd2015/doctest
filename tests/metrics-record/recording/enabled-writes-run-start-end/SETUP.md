# Scenario

**Feature**: metrics-on suite writes run_start and run_end JSONL events

```
# metrics opt-in
RunTest(1-leaf fixture, MetricsOn=true, MetricsRoot=tmp)
  -> one *.jsonl with run_start then run_end (schema_version 1)
```

## Preconditions

- MetricsRoot empty before run (temp dir).
- Default suite (no labels).

## Steps

1. Run suite with metrics enabled.
2. Decode first new JSONL file.
3. Assert presence of run_start and run_end with schema_version 1.

## Context

- run_start should carry mode.default_suite true (or equivalent nested mode object);
  assert type markers first; soft-check default_suite when present.
- run_end includes wall_ns / totals fields when implementer follows P2 schema.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.MetricsOn = true
	return nil
}
```
