# Scenario

**Feature**: `--no-leaf-cache` disables leaf-cache skip after a warm hit

```
run1: test fixture                  -> store pass
run2: test fixture                  -> Cached > 0
run3: test fixture --no-leaf-cache  -> 0 Cached
```

## Preconditions

- Parent pass fixture.
- Flag is leaf-cache-specific (does not require -count).

## Steps

1. Set Args3 to `test <fixture> --no-leaf-cache`.
2. Assert run2 Cached > 0; run3 Cached == 0.

## Context

- Opt-out for users who want pass re-execution without `-count` side effects on go test.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	t.Helper()
	req.Args3 = []string{"test", req.FixtureDir, "--no-leaf-cache"}
	return nil
}
```
