# Scenario

**Feature**: disjoint user package key and vendor go.mod key both present

```
user: project-source -> "from-user"
vendor: active-vendor-gomod -> placeholder
  -> final has both mappings
```

## Preconditions

- No key collision between user package path and vendor go.mod path.

## Steps

1. Seed user package key.
2. Vendor adds go.mod → placeholder via `$ACTIVE_BRIDGE` value sentinel.

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.UserReplace = map[string]string{
		"project-source": "from-user",
	}
	req.VendorReplace = map[string]string{
		"active-vendor-gomod": "$ACTIVE_BRIDGE",
	}
	req.CreateActiveBridgeFile = true
	return nil
}
```
