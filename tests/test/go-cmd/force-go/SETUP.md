# Scenario

**Feature**: `--go-cmd=go` always selects `go` (detection ignored)

```
ResolveGoTestCmd("go", needsXgo=true) -> "go"
# mocks may panic at runtime; CLI still honors force-go
```

## Preconditions

- Force-go leaves may seed a mock-positive graph to prove override.
- No PATH ensure required for `go`.

## Steps

1. Set `GoCmdFlag=go`.
2. Leaf enables detect with transitive mock fixture.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.GoCmdFlag = "go"
	req.ParseOnly = false
	req.CheckAvailable = false
	return nil
}
```
