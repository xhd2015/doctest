# Scenario

**Feature**: source-input change forces gomod-src cache miss and rebuild

```
first write -> seed tidy-done
change exactly one source factor (go.mod | go.sum | modules.txt | assert path)
second WriteGoModWithVendorBridges
  -> rebuild; tidy-done dropped when gen mod/sum wrote
```

## Preconditions

- First write established fingerprint + bridges.
- Leaves change **one** invalidation factor only (MECE).

## Steps

1. Prepare parent module; first write; seed tidy-done.
2. Leaf sets the single change field used by Mode `write-second`.

```go
import (
	"testing"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Mode = "write-second"
	req.ModPath = defaultModPath
	req.HasMod = true
	// Leaves override fixtures (vendor, assert flags, ChangeSource*).
	return nil
}
```
