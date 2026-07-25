# Scenario

**Feature**: --no-color suppresses ANSI in gen-plan result output

```
doctest test --no-color …
  -> gen-plan: result plain text (no ESC on gen-plan lines)
```

## Preconditions

- Same cold fixture as color-on sibling.

## Steps

1. Args with --no-color.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	t.Helper()
	req.ColorMode = "off"
	req.Args = baseTestArgs(req, "--no-color", req.TreeRoot)
	return nil
}
```
