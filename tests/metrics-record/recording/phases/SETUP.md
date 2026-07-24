# Scenario

**Feature**: metrics-on suite emits pipeline phase events

```
# phase spans under MetricsRoot
RunTest(1-leaf fixture, MetricsOn=true)
  -> JSONL includes type=phase for discover, materialize, generate, go_test
```

## Preconditions

- MetricsRoot is an empty temp dir.
- One-leaf pass fixture.

## Steps

1. Run suite with metrics enabled.
2. Decode JSONL events.
3. Assert phase events exist with positive elapsed for core phases.

## Context

- Phases are tree-scoped completed spans (not start/end pairs).
- `go_test` covers compile+execute; discover/materialize/generate are doctest-side.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Op = "record_run"
	req.MetricsOn = true
	return nil
}
```
