# Scenario

**Feature**: warm hit returns the same BridgeRoot paths from bridges cache

```
first bridges[0].BridgeRoot = P
second WriteGoModWithVendorBridges
  -> bridges[0].BridgeRoot == P
```

## Steps

1. Inherit vendor-bridges Setup (first write + SnapBridgeRoot).
2. Confirm snapshot before measured warm call.

```go
import (
	"testing"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	if req.Mode != "write-second" {
		t.Fatalf("warm-returns-cached expects Mode write-second, got %q", req.Mode)
	}
	if req.SnapBridgeRoot == "" {
		t.Fatal("parent must snapshot BridgeRoot from first write")
	}
	if req.SnapBridgeCount < 1 {
		t.Fatal("parent must record non-empty first bridges")
	}
	return nil
}
```
