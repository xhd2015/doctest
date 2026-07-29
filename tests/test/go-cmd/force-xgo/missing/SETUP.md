# Scenario

**Feature**: force xgo but binary missing on search PATH → clear error

```
ResolveGoTestCmd("xgo", _) -> "xgo"
EnsureGoTestCmdAvailable("xgo", fakePATHWithoutXgo) -> error
  message: xgo not found in PATH (actionable)
```

## Preconditions

- `SearchPATH` is a temp directory containing a decoy `go` but **no** `xgo`.
- Lookup must use the injected search PATH only (not fall back to process PATH
  where a real xgo might exist).
- No product PATH shim that renames `go`.

## Steps

1. Force xgo.
2. Enable `CheckAvailable` with `fakePATHWithoutXgo`.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.GoCmdFlag = "xgo"
	req.Detect = false
	req.NeedsXgo = false
	req.CheckAvailable = true
	req.SearchPATH = fakePATHWithoutXgo(t)
	return nil
}
```
