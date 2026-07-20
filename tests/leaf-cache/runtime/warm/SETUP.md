# Scenario

**Feature**: successful leaves are stored and skipped on a warm second run

```
# cold
doctest test fixture -> run leaf -> pass -> PutPass
# warm
doctest test fixture -> GetPass hit -> skip -> Cached > 0
```

## Preconditions

- 1-leaf always-pass fixture.
- Both runs use default flags (leaf-cache enabled; no -count/force).

## Steps

1. Write pass fixture.
2. Args = Args2 = `test <fixture>` (no disable flags).

## Context

- Baseline happy path for P2 Cached reporting.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	t.Helper()
	req.FixtureDir = preparePassFixture(t, 1)
	req.Args = []string{"test", req.FixtureDir}
	req.Args2 = []string{"test", req.FixtureDir}
	return nil
}
```
