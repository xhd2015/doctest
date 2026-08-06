# Scenario

**Feature**: user seed key and hook key both appear when disjoint

```
user seed: project-source -> "from-user"
hook write: active-vendor -> "from-hook"
  -> final has both keys
```

## Preconditions

- Different disk-path keys; no same-key conflict.

## Steps

1. Seed user on `project-source`.
2. Hook writes `active-vendor` only.

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.UserReplace = map[string]string{
		"project-source": "from-user",
	}
	req.HookOverlays = [][]OverlayEntry{
		{{Source: "active-vendor", Replace: "from-hook"}},
	}
	return nil
}
```
