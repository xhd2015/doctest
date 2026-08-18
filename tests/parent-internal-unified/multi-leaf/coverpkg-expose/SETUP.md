# Scenario

**Feature**: parent-internal expose + CI-style `-coverpkg=mod/...` leaves a
coverprofile that downstream `go tool cover` can consume (no session-generated
expose packages in the final profile)

```
RunTest(multi-leaf parent internal,
  Cover + CoverProfile + CoverPkg=example.com/app/... + CoverMode=set)
  -> exit 0 (suite package)
  -> cover.out non-empty
  -> no go tool cover open failure on __doctest_internal_expose/*/expose.go during run
  -> cover.out has no lines for session-generated expose facades
  -> go tool cover -func=cover.out from ModuleRoot succeeds
```

Crime scene (scaff CI / reconstruct 2026-08-18): with product-module coverpkg,
doctest PASS writes expose paths into `-coverprofile`, then cleans product
expose files → `go tool cover -func` / scaff merge report fails:
`no required module provides package …/__doctest_internal_expose/…`.

## Preconditions

- Fixture from parent multi-leaf Setup (imports `internal/greet` → expose).
- Cover profile path absolute under `t.TempDir()` (parallel-safe).
- `CoverPkg` uses product module path `example.com/app/...` (same wildcard shape
  as CI `github.com/<mod>/...`).

## Steps

1. Enable cover + absolute CoverPath.
2. Set CoverPkg + CoverMode to match CI recipe.
3. Run subject tree; Assert success, clean profile, and `go tool cover -func`.

```go
import (
	"os"
	"path/filepath"
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	t.Helper()
	_ = d
	req.WithCover = true
	req.CoverPath = filepath.Join(t.TempDir(), "cover.out")
	_ = os.Remove(req.CoverPath)
	// Product module wildcard — instruments internal/greet and expose.
	req.CoverPkg = modPath + "/..."
	req.CoverMode = "set"
	return nil
}
```
