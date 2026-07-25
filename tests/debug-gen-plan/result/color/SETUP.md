# Scenario

**Feature**: result hierarchy honors --color / --no-color for status SGR

```
--color:
  gray=unchanged  green=modified  red=deleted  (CSI on path lines)
--no-color:
  no ESC sequences on gen-plan result lines
```

## Preconditions

- Single-tree fixture; cold gen so modified paths exist for green.

## Steps

1. prepareSingleTree; leaves set color flag in Args.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	t.Helper()
	prepareSingleTree(t, req)
	return nil
}
```
