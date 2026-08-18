# Scenario

**Feature**: parent-internal expose must tolerate CI-style `-coverpkg=mod/...`

```
RunTest(multi-leaf parent internal,
  Cover + CoverProfile + CoverPkg=example.com/app/... + CoverMode=set)
  -> exit 0 (suite package)
  -> cover.out non-empty
  -> no go tool cover open failure on __doctest_internal_expose/*/expose.go
```

Crime scene (scaff CI / reconstruct): with product-module coverpkg, `go tool cover`
opens the logical expose path that exists only via `-overlay` →
`no such file or directory` → suite build failed.

## Preconditions

- Fixture from parent multi-leaf Setup (imports `internal/greet` → expose).
- Cover profile path absolute under `t.TempDir()` (parallel-safe).
- `CoverPkg` uses product module path `example.com/app/...` (same wildcard shape
  as CI `github.com/<mod>/...`).

## Steps

1. Enable cover + absolute CoverPath.
2. Set CoverPkg + CoverMode to match CI recipe.
3. Run subject tree; Assert success and profile.

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
