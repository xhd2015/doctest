# Scenario

**Feature**: --color paints result statuses with green and/or gray SGR

```
doctest test --color … 
  -> gen-plan: result lines include CSI (green modified, gray unchanged)
```

## Preconditions

- Cold gen so at least one modified (green) path is expected.

## Steps

1. Args with --color.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	t.Helper()
	req.ColorMode = "on"
	req.Args = baseTestArgs(req, "--color", req.TreeRoot)
	return nil
}
```
