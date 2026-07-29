# Scenario

**Feature**: `--go-cmd` default **auto** chooses go vs xgo from mock detection

```
# mode auto (flag omitted or --go-cmd=auto)
DetectXgoMockUsage(entrypoints)
  -> needsXgo=false -> ResolveGoTestCmd -> "go"
  -> needsXgo=true  -> ResolveGoTestCmd -> "xgo"
```

## Preconditions

- Leaves under `auto/` use detection against a temp module fixture.
- No forced mode: `GoCmdFlag` empty or `auto`.
- Do not require a real `xgo` binary for resolve-only leaves.

## Steps

1. Set mode to auto (omitted flag semantics via empty `GoCmdFlag`).
2. Leaf seeds no-mock or transitive-mock module and runs detect+resolve.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	// Empty GoCmdFlag == omitted / auto (product treats "" like auto).
	req.GoCmdFlag = ""
	req.Detect = true
	req.CheckAvailable = false
	req.ParseOnly = false
	return nil
}
```
