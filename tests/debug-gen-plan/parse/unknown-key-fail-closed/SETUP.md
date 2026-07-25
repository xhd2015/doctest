# Scenario

**Feature**: unknown DOCTEST_DEBUG keys still fail closed after gen-plan lands

```
debug.Parse("bypass-go-test=1,not-a-key=1")
  -> error containing unknown key
```

## Preconditions

- Uses only already-known + unknown keys so this leaf can stay GREEN once
  gen-plan is added (does not depend on gen-plan itself).

## Steps

1. Set DebugEnv to a known key plus an unknown key.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	t.Helper()
	// Avoid gen-plan here so pre-feature and post-feature both fail on not-a-key.
	req.DebugEnv = "bypass-go-test=1,not-a-key=1"
	return nil
}
```
