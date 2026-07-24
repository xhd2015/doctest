# Scenario

**Feature**: warm leaf-cache skips emit grey progress dots under quiet + --color

```
run1: test 2-pass fixture -> store both (cold; color optional)
run2: test --color same fixture -> Cached >= 2; progress has grey dots
```

## Preconditions

- Two-leaf always-pass fixture (`preparePassFixture(t, 2)`).
- Run2 forces color with `--color` (pipe-safe).
- Leaf-cache enabled (no `-count` / `-a` / `--no-leaf-cache` on run2).

## Steps

1. Prepare 2-pass fixture.
2. Args = `test <fixture>`; Args2 = `test --color <fixture>`.
3. Assert run2 Cached >= 2 and grey progress-dot count >= 2.

## Context

- Warm skip identities are in the this-run skip set; progress `.` for those
  leaves must be `\x1b[90m.\x1b[0m`, not plain `.`.
- Executed-pass plain dots are covered indirectly (run1 may be plain); fail red
  is a sibling leaf.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	t.Helper()
	req.Op = "runtime_multi"
	req.FixtureDir = preparePassFixture(t, 2)
	req.Args = []string{"test", req.FixtureDir}
	req.Args2 = []string{"test", "--color", req.FixtureDir}
	return nil
}
```
