# Scenario

**Feature**: on a shared go.mod Replace key, vendor-gomod mapping wins over user seed

```
user seed: active-vendor-gomod -> "from-user-gomod"
           project-source      -> "seed-only"
vendor:    active-vendor-gomod -> placeholder go.mod path
  -> final[go.mod] == placeholder (vendor wins)
  -> final[project-source] == "seed-only" (proves user seed applied)
```

## Preconditions

- Same disk-path key (project vendor go.mod path).
- Vendor value is the fixture placeholder path (`ActiveBridgeSource` on Response).
- Extra seed-only key proves user layer ran.

## Steps

1. Seed user map on `active-vendor-gomod` + seed-only package key.
2. VendorReplace maps go.mod key to `$ACTIVE_BRIDGE` (Run expands to placeholder).

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.UserReplace = map[string]string{
		"active-vendor-gomod": "from-user-gomod",
		"project-source":      "seed-only",
	}
	// Vendor later-wins: same key, value is placeholder (sentinel expanded in Run).
	req.VendorReplace = map[string]string{
		"active-vendor-gomod": "$ACTIVE_BRIDGE",
	}
	req.CreateActiveBridgeFile = true
	return nil
}
```
