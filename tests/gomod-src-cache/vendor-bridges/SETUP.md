# Scenario

**Feature**: vendor bridges cache returns mappings on warm hit without rewriting placeholders

```
parent with nogo vendor module (no go.mod)
first WriteGoModWithVendorBridges -> bridges + overlay placeholder
second identical write
  -> same BridgeRoot returned from doctest.vendor-bridges.json
  -> placeholder mtime stable
```

## Preconditions

- Vendor module lacks go.mod so a non-empty bridges list is produced.
- Layout remains replace → project vendor + overlay placeholder (no hardlink bridge).

## Steps

1. Prepare parent + seedVendorNogo; first write; snapshot bridges/placeholder.
2. Mode write-second measured call with no source mutations.

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
		t.Fatal("expected non-empty bridges for nogo vendor module")
	}
	snapshotBridges(t, req, bridges)
	seedTidyDone(t, req.GenDir)
	return nil
}
```
