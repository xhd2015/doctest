# Scenario

**Feature**: user overlay alone materializes one `-overlay=` with user Replace

```
MaterializeUserVendorOverlay(user.json, "", dest)
  -> GoFlags == ["-overlay=<file>"]
  -> Replace == user map
```

## Preconditions

- No pre_test hooks; no vendor overlay.
- User Replace has one package-path key (stable string paths under fixture).

## Steps

1. Seed `UserReplace` with a synthetic absolute key/value pair.
2. Force materialize helper path (`UseMaterializeHelper`).

```go
import (
	"path/filepath"
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	// Stable synthetic paths (not required to exist on disk for overlay JSON).
	src := filepath.Join(string(filepath.Separator), "user-only-src", "pkg", "a.go")
	req.UserReplace = map[string]string{src: "/replacement/user-only.go"}
	req.UseMaterializeHelper = true
	req.PreTest = nil
	return nil
}
```
