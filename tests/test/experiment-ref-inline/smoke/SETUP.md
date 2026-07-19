# Scenario

**Feature**: classic generation still runs for a tiny tree (P0 regression smoke)

```
# mini suite run with Options.ExperimentRefInsteadOfInline set or clear
RunTest(1-leaf pass fixture, ExperimentRefInsteadOfInline=?)
  -> classic AssembleTestSource (P0) -> leaf pass -> no error
```

## Preconditions

- Uses package suite entry (`runner.RunTest`).
- Fixture is a one-leaf pass tree from `testtree.WritePassFailTree`.
- Does **not** inspect generated source for refs (P1).

## Steps

1. Set `Op=mini_run` and the Options bool.
2. Assert `RunErr` empty (suite succeeds).

## Context

- Proves the default path and flag-on path do not break classic gen in P0.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Op = "mini_run"
	return nil
}
```
