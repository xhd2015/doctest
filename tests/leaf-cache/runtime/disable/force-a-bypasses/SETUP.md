# Scenario

**Feature**: `-a` force disables leaf-cache skip after a warm hit

```
run1: test fixture     -> store pass
run2: test fixture     -> Cached > 0
run3: test fixture -a  -> 0 Cached
```

## Preconditions

- Parent pass fixture; `-a` is the short force flag.

## Steps

1. Set Args3 to `test <fixture> -a`.
2. Assert run2 Cached > 0; run3 Cached == 0.

## Context

- Force re-runs leaves even when the pass store would hit.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	t.Helper()
	req.Args3 = []string{"test", req.FixtureDir, "-a"}
	return nil
}
```
