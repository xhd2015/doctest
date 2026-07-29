# Scenario

**Feature**: `--go-cmd=go` forces `go` even when transitive mock is detected

```
DetectXgoMockUsage(runpkg→helper→mock) -> needsXgo=true
ResolveGoTestCmd("go", true) -> "go"   # force wins
```

## Preconditions

- Same transitive mock fixture as `auto/transitive-mock`.
- Mode is forced `go`, not auto.

## Steps

1. Seed transitive mock module.
2. Force `GoCmdFlag=go`; detect + resolve.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	seedTransitiveMockModule(t, req)
	req.GoCmdFlag = "go"
	req.Detect = true
	return nil
}
```
