# Scenario

**Feature**: second WriteGoModWithVendorBridges with identical sources is a warm hit

```
first write -> seed tidy-done -> force old go.mod mtime
second write (same sources/flags)
  -> go.mod mtime stable
  -> tidy-done retained
  -> gomod-src content stable
```

## Preconditions

- First write already populated gen root cache artifacts.
- Desired source inputs unchanged between calls.

## Steps

1. Prepare fresh dirs; first WriteGoModWithVendorBridges in Setup.
2. Seed tidy-done; snapshot go.mod mtime and gomod-src content.
3. Mode `write-second` is the measured warm call (no source mutations).

```go
import (
	"testing"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Mode = "write-second"
	req.ModPath = defaultModPath
	req.HasMod = true
	prepareFreshGen(t, req, "module example.com/app\n\ngo 1.21\n")
	firstWrite(t, req)
	seedTidyDone(t, req.GenDir)
	snapshotGoModMtime(t, req)
	snapshotGomodSrc(t, req)
	return nil
}
```
