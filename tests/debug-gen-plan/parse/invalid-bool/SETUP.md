# Scenario

**Feature**: gen-plan rejects non-bool values after the key is known

```
debug.Parse("gen-plan=maybe")
  -> error (invalid bool), not silent accept
```

## Preconditions

- Classic TDD: RED until gen-plan is recognized (pre-feature error is "unknown
  key"; post-feature must be bool validation error).

## Steps

1. Set DebugEnv to `gen-plan=maybe`.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	t.Helper()
	req.DebugEnv = "gen-plan=maybe"
	return nil
}
```
