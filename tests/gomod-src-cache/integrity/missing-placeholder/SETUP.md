# Scenario

**Feature**: missing overlay placeholder listed in bridges forces integrity miss rebuild

```
first write with nogo vendor module (creates placeholder)
delete placeholder path from bridges
second write (same sources)
  -> rebuild; placeholder restored
```

## Steps

1. Parent go.mod + seedVendorNogo; first write captures BridgeRoot.
2. Set DeletePlaceholder so Run removes placeholder before second call.

```go
import (
	"testing"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Mode = "write-second"
	req.ModPath = defaultModPath
	req.HasMod = true
	prepareFreshGen(t, req, "module example.com/app\n\ngo 1.21\n")
	seedVendorNogo(t, req)
	bridges := firstWrite(t, req)
	if len(bridges) == 0 {
		t.Fatal("expected at least one vendor bridge for nogo module")
	}
	snapshotBridges(t, req, bridges)
	if req.SnapPlaceholderPath == "" {
		t.Fatal("expected placeholder path from first bridges")
	}
	if !fileExists(req.SnapPlaceholderPath) {
		t.Fatalf("placeholder missing after first write: %s", req.SnapPlaceholderPath)
	}
	req.DeletePlaceholder = true
	return nil
}
```
