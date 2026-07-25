# Scenario

**Feature**: gen-plan result phase reprints hierarchy with statuses after generate

```
after generate [+ prune]:
  gen-plan: result
    <same hierarchy as plan>
    # colors when on: gray unchanged, green modified, red deleted
    summary: modified=N unchanged=M deleted=K
```

## Preconditions

- Product CLI; DebugEnv includes gen-plan + bypass-go-test.
- Single-tree fixtures for status/color leaves (same hierarchy as plan).

## Steps

1. Group defaults Mode=cli and DebugEnv combo.
2. Status/color leaves prepare fixtures and color flags.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	t.Helper()
	req.Mode = "cli"
	req.DebugEnv = debugGenPlanBypass
	return nil
}
```
