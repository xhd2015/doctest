# Scenario

**Feature**: gen-plan pairs with bypass-go-test in one DOCTEST_DEBUG string

```
debug.Parse("gen-plan=1,bypass-go-test=1")
  -> GenPlan true, BypassGoTest true
```

## Preconditions

- Classic TDD: RED until gen-plan is known (bypass alone already works).

## Steps

1. Set DebugEnv to combined keys.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	t.Helper()
	req.DebugEnv = "gen-plan=1,bypass-go-test=1"
	return nil
}
```
