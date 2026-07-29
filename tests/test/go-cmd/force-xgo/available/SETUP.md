# Scenario

**Feature**: `--go-cmd=xgo` resolves to `xgo` even when no mock is present

```
ResolveGoTestCmd("xgo", needsXgo=false) -> "xgo"
# no EnsureGoTestCmdAvailable (avoid requiring real xgo binary)
```

## Preconditions

- Injected `NeedsXgo=false` (or no detect) proves force ignores detection.
- Do not call PATH ensure in this leaf.

## Steps

1. Force xgo with needsXgo=false.
2. Resolve only.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.GoCmdFlag = "xgo"
	req.Detect = false
	req.NeedsXgo = false
	req.CheckAvailable = false
	return nil
}
```
