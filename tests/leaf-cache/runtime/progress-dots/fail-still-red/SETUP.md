# Scenario

**Feature**: fail progress dots stay red when color is on (regression lock)

```
doctest test --color pass+fail fixture
  -> fail leaf -> progress <ansi red .>
  -> pass leaf -> plain . (not required green)
```

## Preconditions

- 1-pass + 1-fail fixture via `preparePassFailMixFixture`.
- `--color` on the single nested run.

## Steps

1. Prepare mix fixture.
2. Args = `test --color <fixture>` (no Args2).
3. Assert non-zero exit and >= 1 red progress dot.

## Context

- Complements grey warm dots; ensures fail coloring is not broken when greying
  cached skips. Product already reds fail dots — this leaf is a lock-in.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	t.Helper()
	req.Op = "runtime_multi"
	req.FixtureDir = preparePassFailMixFixture(t, 1, 1)
	// Single invocation: leave Args2 empty so Run only executes Args once.
	// (runtime_multi defaults Args2=Args when empty — set Args2 to a no-op
	// re-run is fine; fail twice is OK for red-dot presence.)
	req.Args = []string{"test", "--color", req.FixtureDir}
	req.Args2 = []string{"test", "--color", req.FixtureDir}
	return nil
}
```
