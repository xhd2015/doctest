# Scenario

**Feature**: Parse accepts gen-plan=1 and sets GenPlan

```
debug.Parse("gen-plan=1")
  -> err == nil
  -> Settings.GenPlan == true
```

## Preconditions

- Classic TDD: RED until `gen-plan` is a known key and Settings.GenPlan exists.

## Steps

1. Set DebugEnv to `gen-plan=1`.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	t.Helper()
	req.DebugEnv = "gen-plan=1"
	return nil
}
```
