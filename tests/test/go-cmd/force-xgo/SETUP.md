# Scenario

**Feature**: `--go-cmd=xgo` always selects `xgo` (detection ignored)

```
ResolveGoTestCmd("xgo", needsXgo=false|true) -> "xgo"
EnsureGoTestCmdAvailable("xgo", PATH) -> ok | not found
```

## Preconditions

- Force mode does not depend on mock detection for the resolve result.
- `available/` skips PATH ensure (no real xgo required).
- `missing/` injects a fake search PATH with no xgo binary.

## Steps

1. Set `GoCmdFlag=xgo`.
2. Leaves set detect/ensure variants.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.GoCmdFlag = "xgo"
	req.ParseOnly = false
	// Detection optional under force; default off unless leaf enables.
	req.Detect = false
	req.NeedsXgo = false
	return nil
}
```
