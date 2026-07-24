# Scenario

**Feature**: leaf-cache skip is disabled by count / force / no-leaf-cache flags

```
# populate store
run1: doctest test fixture (default) -> PutPass
# prove warm skip works
run2: doctest test fixture (default) -> Cached > 0
# disable skip
run3: doctest test fixture <disable-flag> -> 0 Cached
```

## Preconditions

- 1-pass fixture; leaf-cache enabled for run1/run2.
- Children set Args3 with exactly one disable mechanism (MECE).
- Three-run sequence avoids vacuous "0 Cached" passes before skip is implemented.

## Steps

1. Prepare pass fixture; Args = Args2 = default `test <fixture>`.
2. Child sets Args3 with the disable flag under test.
3. Assert run2 Cached > 0 and run3 Cached == 0.

## Context

- Sibling leaves are MECE disable knobs: `-count`, `-a`, `--no-leaf-cache`.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	t.Helper()
	req.FixtureDir = preparePassFixture(t, 1)
	req.Args = []string{"test", req.FixtureDir}
	req.Args2 = []string{"test", req.FixtureDir}
	// Args3 set by child with disable flag
	return nil
}
```
