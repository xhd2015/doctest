# Scenario

**Feature**: `-count=1` disables leaf-cache skip so quiet+color emits no grey warm-skip dots

```
run1: test 2-pass -> store
run2: test --color -> warm grey dots + Cached >= 2 (precondition)
run3: test --color -count=1 -> 0 Cached; 0 grey progress dots
```

## Preconditions

- Two-leaf always-pass fixture.
- Run2 proves warm grey path is active before testing `-count` bypass.
- Run3 uses both `--color` and `-count=1`.

## Steps

1. Prepare 2-pass fixture.
2. Args store; Args2 warm+color; Args3 count bypass+color.
3. Assert run3 Cached == 0 and grey progress dots == 0.

## Context

- Any `-count` disables programmatic skip; executed passes stay plain `.`.
- Locks that grey dots are tied to leaf-cache skip set, not merely "second run".

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	t.Helper()
	req.Op = "runtime_multi"
	req.FixtureDir = preparePassFixture(t, 2)
	req.Args = []string{"test", req.FixtureDir}
	req.Args2 = []string{"test", "--color", req.FixtureDir}
	req.Args3 = []string{"test", "--color", "-count=1", req.FixtureDir}
	return nil
}
```
