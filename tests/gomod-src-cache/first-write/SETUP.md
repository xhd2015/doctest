# Scenario

**Feature**: cold gen root first WriteGoModWithVendorBridges creates cache artifacts

```
empty genDir + parent go.mod (no vendor)
  -> WriteGoModWithVendorBridges
  -> doctest.gomod-src + doctest.vendor-bridges.json
```

## Preconditions

- Gen root starts empty (no prior fingerprint or bridges cache).
- Parent module has go.mod; vendor is intentionally absent for baseline cold write.

## Steps

1. Prepare fresh ModRoot/GenDir with a minimal parent go.mod.
2. Mode `write-once` is the measured cold write.

```go
import (
	"testing"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Mode = "write-once"
	req.ModPath = defaultModPath
	req.HasMod = true
	prepareFreshGen(t, req, "module example.com/app\n\ngo 1.21\n")
	return nil
}
```
